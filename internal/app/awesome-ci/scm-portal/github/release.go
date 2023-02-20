package github

import (
	"errors"
	"fmt"
	"time"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
)

// CreateRelease
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

// PublishRelease
func (ghrc *GitHubRichClient) PublishRelease(
	tagName string,
	releasePrefix string,
	releaseBranch string,
	body string,
	releaseId int64,
	uploadArtifacts []tools.UploadAsset) (releaseAssets []*github.ReleaseAsset, err error) {

	if releaseId == 0 {
		log.Infoln("no release found, creating one...")
		release, err := ghrc.CreateRelease(tagName, releasePrefix, releaseBranch, body)
		if err != nil {
			log.Fatalln(err)
		}
		releaseId = *release.ID
	}

	existingRelease, _, err := ghrc.Client.Repositories.GetRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseId)
	if err != nil {
		return
	}

	// upload any given artifacts
	var releaseBodyAssets string = ""
	if len(uploadArtifacts) > 0 {
		releaseBodyAssets = "### Asstes\n"

		for _, fileAndInfo := range uploadArtifacts {
			log.Infof("uploading %s as asset to release\n", fileAndInfo.Name)
			// Upload assets to GitHub Release
			relAsset, _, err := ghrc.Client.Repositories.UploadReleaseAsset(
				ctx,
				ghrc.Owner,
				ghrc.Repository,
				releaseId,
				&github.UploadOptions{
					Name: fileAndInfo.Name,
				},
				&fileAndInfo.File)
			if err != nil {
				log.Println("error at uploading asset to release: ", err)
			} else {
				// add asset to release body
				releaseBodyAssets = fmt.Sprintf("%s\n- [%s](%s) `%s`\n  Sha256: `%x`", releaseBodyAssets, fileAndInfo.Name, *relAsset.BrowserDownloadURL, fileAndInfo.Infos.ModTime().Format(time.RFC3339), fileAndInfo.Hash)

				releaseAssets = append(releaseAssets, relAsset)
			}
		}
	}

	newReleaseBody := fmt.Sprintf("%s\n\n%s", *existingRelease.Body, releaseBodyAssets)
	existingRelease.Body = &newReleaseBody

	// publishing release
	*existingRelease.Draft = false
	_, _, err = ghrc.Client.Repositories.EditRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseId,
		existingRelease)
	if err != nil {
		return releaseAssets, err
	}

	return
}

// GetLatestReleaseVersion
func (ghrc *GitHubRichClient) GetLatestReleaseVersion() (latestRelease *github.RepositoryRelease, err error) {
	var releaseMap = make(map[string]*github.RepositoryRelease)
	var loadReleasesFromRepo func(page int)

	loadReleasesFromRepo = func(page int) {
		log.Tracef("querying release page %d", page)
		releases, response, err := ghrc.Client.Repositories.ListReleases(ctx, ghrc.Owner, ghrc.Repository, &github.ListOptions{
			PerPage: 100,
			Page:    page,
		})
		if err != nil {
			log.Fatalf("error at loading repos from emst: %v", err)
		}

		log.Tracef("found %d releases at page %d begin with mapping...", len(releases), page)
		for _, release := range releases {
			releaseMap[*release.TagName] = release
		}

		log.Traceln("####### next page:", response.NextPage)
		if response.NextPage == 0 {
			log.Traceln("ended with paging through releases")
			return
		} else {
			loadReleasesFromRepo(page + 1)
		}
	}

	// starting recusive function
	loadReleasesFromRepo(1)

	return ghrc.findLatestRelease(`.`, releaseMap)
}

func (ghrc *GitHubRichClient) findLatestRelease(directory string, githubReleaseMap map[string]*github.RepositoryRelease) (latestRelease *github.RepositoryRelease, err error) {
	log.Traceln("open local git repository...")
	gitRepo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}
	log.Traceln("...opened local git repository")

	tags, err := tools.GetGitTagsUpToHead(gitRepo)

	log.Traceln(tags)
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
