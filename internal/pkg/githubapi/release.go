package githubapi

import (
	"awesome-ci/internal/pkg/tools"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v44/github"
)

// CreateRelease
func (ghrc *GitHubRichClient) CreateRelease(version string, releaseBranch string, body string, draft bool) (createdRelease *github.RepositoryRelease, err error) {

	relName := "Release " + version
	if releaseBranch == "" {
		releaseBranch = tools.GetDefaultBranch()
	}

	// get body for release
	if body != "" {
		bodyFile, err := tools.CheckIsFile(body)
		if err == nil {
			body = bodyFile
		}
	}

	releaseObject := github.RepositoryRelease{
		TargetCommitish: &releaseBranch,
		TagName:         &version,
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

	envVars, err := OpenEnvFile()
	if err != nil {
		return nil, err
	}
	envVars.Set("ACI_RELEASE_ID", fmt.Sprintf("%d", *createdRelease.ID))
	err = envVars.SaveEnvFile()
	return
}

// PublishRelease
func (ghrc *GitHubRichClient) PublishRelease(version string, releaseBranch string, body string, releaseId int64, uploadArtifacts *string) (err error) {
	draftFalse := false

	if releaseId == 0 {
		releaseIdStr, releaseIdBool := os.LookupEnv("ACI_RELEASE_ID")
		if !releaseIdBool {
			log.Println("No release found, creating one...")
			release, err := ghrc.CreateRelease(version, releaseBranch, body, true)
			if err != nil {
				log.Fatalln(err)
			}
			releaseId = *release.ID
		} else {
			releaseId, err = strconv.ParseInt(releaseIdStr, 10, 64)
			if err != nil {
				fmt.Printf("%s of type %T", releaseIdStr, releaseIdStr)
				os.Exit(2)
			}
		}
	}

	existingRelease, _, err := ghrc.Client.Repositories.GetRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseId)
	if err != nil {
		return err
	}

	// upload any given artifacts
	var releaseBodyAssets string = ""
	if *uploadArtifacts != "" {
		filesAndInfos, err := tools.GetAsstes(uploadArtifacts, false)
		if err != nil {
			return err
		}

		releaseBodyAssets = "### Asstes\n"

		for i, fileAndInfo := range filesAndInfos {
			fmt.Printf("uploading %s as asset to release\n", fileAndInfo.Name)
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

				// export Download URL to env. See: #53
				envVars, err := OpenEnvFile()
				if err != nil {
					log.Println("could open envs:", err)
				}
				envVars.Set(fmt.Sprintf("ACI_ARTIFACT_%d_URL", i+1), *relAsset.BrowserDownloadURL)
				err = envVars.SaveEnvFile()
				if err != nil {
					log.Println("could not export atrifact url:", err)
				}
			}
		}
	}

	newReleaseBody := fmt.Sprintf("%s\n\n%s", *existingRelease.Body, releaseBodyAssets)
	existingRelease.Body = &newReleaseBody

	// publishing release
	*existingRelease.Draft = draftFalse
	_, _, err = ghrc.Client.Repositories.EditRelease(
		ctx,
		ghrc.Owner,
		ghrc.Repository,
		releaseId,
		existingRelease)
	if err != nil {
		return err
	}

	return
}

// GetLatestReleaseVersion
func (ghrc *GitHubRichClient) GetLatestReleaseVersion() (latestRelease *github.RepositoryRelease, err error) {
	var releaseMap = make(map[string]*github.RepositoryRelease)

	var listOptions github.ListOptions = github.ListOptions{
		PerPage: 100,
		Page:    0,
	}

	for listOptions.Page >= 0 {
		releases, response, err := ghrc.Client.Repositories.ListReleases(ctx, ghrc.Owner, ghrc.Repository, &listOptions)

		if err != nil {
			return nil, err
		}

		for _, release := range releases {
			releaseMap[*release.TagName] = release
		}

		if listOptions.Page == response.NextPage {
			break
		}

		listOptions.Page = response.NextPage
	}

	return ghrc.findLatestRelease(`.`, releaseMap)
}

func (ghrc *GitHubRichClient) findLatestRelease(directory string, githubReleaseMap map[string]*github.RepositoryRelease) (latestRelease *github.RepositoryRelease, err error) {
	gitRepo, err := git.PlainOpen(directory)

	if err != nil {
		return nil, err
	}

	tags, err := tools.GetGitTagsUpToHead(gitRepo)

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
