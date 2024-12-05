package github

import (
	"errors"
	"fmt"
	"time"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v67/github"
	"github.com/sirupsen/logrus"
)

// CreateRelease creates a new release in the GitHub repository.
//
// Parameters:
// - tagName: the name of the tag for the release.
// - releasePrefix: the prefix for the release name.
// - releaseBranch: the branch for the release (default: git default).
// - body: the body of the release (markdown string or file).
//
// Returns:
// - createdRelease: the created repository release.
// - err: an error if there was a problem creating the release.
func (ghrc *GitHubRichClient) CreateRelease(tagName string, releasePrefix string, releaseBranch string, body string) (createdRelease *github.RepositoryRelease, err error) {
	draft := true
	relName := fmt.Sprintf("%s %s", releasePrefix, tagName)
	if releaseBranch == "" {
		releaseBranch = tools.GetDefaultBranch()
	}

	releaseObject := github.RepositoryRelease{
		TargetCommitish: &releaseBranch,
		TagName:         &tagName,
		Name:            &relName,
		Draft:           &draft,
		Body:            &body,
	}
	createdRelease, _, err = ghrc.Client.Repositories.CreateRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		&releaseObject)
	if err != nil {
		err = fmt.Errorf("error at creating github release: %v", err)
		return
	}
	return
}

// PublishRelease publishes a GitHub release.
//
// Parameters:
// - tagName: the name of the tag for the release.
// - releasePrefix: the prefix for the release name.
// - releaseBranch: the branch for the release.
// - body: the body of the release.
// - releaseID: the ID of the existing release (0 if not found).
// - uploadArtifacts: the list of artifacts to upload.
//
// Returns:
// - releaseAssets: the list of release assets.
// - err: an error if any occurred.
func (ghrc *GitHubRichClient) PublishRelease(
	tagName string,
	releasePrefix string,
	releaseBranch string,
	body string,
	releaseID int64,
	uploadArtifacts []tools.UploadAsset) (releaseAssets []*github.ReleaseAsset, err error) {

	if releaseID == 0 {
		logrus.Infoln("no release found, creating one...")
		release, err := ghrc.CreateRelease(tagName, releasePrefix, releaseBranch, body)
		if err != nil {
			logrus.Fatalln(err)
		}
		releaseID = *release.ID
	}

	existingRelease, _, err := ghrc.Client.Repositories.GetRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseID)
	if err != nil {
		return
	}

	// upload any given artifacts
	var releaseBodyAssets = ""
	if len(uploadArtifacts) > 0 {
		releaseBodyAssets = "### Assets\n"

		for _, fileAndInfo := range uploadArtifacts {
			logrus.Infof("uploading %s as asset to release\n", fileAndInfo.Name)
			// Upload assets to GitHub Release
			relAsset, _, err := ghrc.Client.Repositories.UploadReleaseAsset(
				ctx,
				ghrc.Owner,
				ghrc.Repository,
				releaseID,
				&github.UploadOptions{
					Name: fileAndInfo.Name,
				},
				&fileAndInfo.File)
			if err != nil {
				logrus.Println("error at uploading asset to release: ", err)
			} else {
				// add asset to release body
				releaseBodyAssets = fmt.Sprintf("%s\n- [%s](%s) `%s`\n  Sha256: `%x`", releaseBodyAssets, fileAndInfo.Name, *relAsset.BrowserDownloadURL, fileAndInfo.Infos.ModTime().Format(time.RFC3339), fileAndInfo.Hash)

				releaseAssets = append(releaseAssets, relAsset)
			}
		}
	}

	var newReleaseBody = ""
	if *existingRelease.Body == "" {
		newReleaseBody = fmt.Sprintf("%s\n\n%s", body, releaseBodyAssets)
	} else {
		newReleaseBody = fmt.Sprintf("%s\n\n%s", *existingRelease.Body, releaseBodyAssets)
	}
	existingRelease.Body = &newReleaseBody

	// publishing release
	*existingRelease.Draft = false
	_, _, err = ghrc.Client.Repositories.EditRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseID,
		existingRelease)
	if err != nil {
		return releaseAssets, err
	}

	return
}

// GetLatestReleaseVersion retrieves the latest release version from the GitHub repository.
//
// It returns the latest non-draft release and any error encountered during the process.
// The latest release is mapped to its tag name in the releaseMap.
// The function queries the releases in the repository in pages of 100, starting from page 1.
// For each release, it checks if it is not a draft and adds it to the releaseMap.
// The function continues querying releases until there are no more pages.
// Finally, it calls the findLatestRelease function to find the latest release based on the releaseMap.
//
// Parameters:
// - ghrc: a pointer to the GitHubRichClient struct.
//
// Returns:
// - latestRelease: a pointer to the latest non-draft release.
// - err: an error if any occurred during the process.
func (ghrc *GitHubRichClient) GetLatestReleaseVersion() (latestRelease *github.RepositoryRelease, err error) {
	var releaseMap = make(map[string]*github.RepositoryRelease)
	var loadReleasesFromRepo func(page int)

	loadReleasesFromRepo = func(page int) {
		logrus.Tracef("querying release page %d", page)
		releases, response, err := ghrc.Client.Repositories.ListReleases(ctx, ghrc.Owner, ghrc.Repository, &github.ListOptions{
			PerPage: 100,
			Page:    page,
		})
		if err != nil {
			logrus.Fatalf("error at loading repos from: %v", err)
		}

		logrus.Tracef("found %d releases at page %d begin with mapping...", len(releases), page)
		for _, release := range releases {
			if !*release.Draft {
				releaseMap[*release.TagName] = release
			} else {
				logrus.Infof("ignoring draft release %s", *release.Name)
			}
		}

		logrus.Traceln("####### next page:", response.NextPage)
		if response.NextPage == 0 {
			logrus.Traceln("ended with paging through releases")
			return
		} else {
			loadReleasesFromRepo(page + 1)
		}
	}

	// starting recusive function
	loadReleasesFromRepo(1)

	return ghrc.findLatestRelease(`.`, releaseMap)
}

// findLatestRelease finds the latest release in the given directory based on the tags in the local git repository.
//
// Parameters:
// - directory: the directory of the local git repository.
// - githubReleaseMap: a map of tags to github.RepositoryRelease objects.
//
// Returns:
// - latestRelease: the latest release found in the directory, or nil if no release is found.
// - err: an error if any occurred during the process.
func (ghrc *GitHubRichClient) findLatestRelease(directory string, githubReleaseMap map[string]*github.RepositoryRelease) (latestRelease *github.RepositoryRelease, err error) {
	logrus.Traceln("open local git repository...")
	gitRepo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}
	logrus.Traceln("...opened local git repository")

	tags, err := tools.GetGitTagsUpToHead(gitRepo)

	logrus.Traceln(tags)
	if err != nil {
		return nil, err
	}

	for i := len(tags) - 1; i >= 0; i-- {
		if latestRelease, found := githubReleaseMap[tags[i].String()]; found {
			return latestRelease, nil
		}
	}

	return nil, errors.New("could not find latest release")
}
