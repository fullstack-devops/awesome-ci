package acigithub

import (
	"awesome-ci/src/models"
	"awesome-ci/src/tools"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/v39/github"
)

// CreateRelease
func CreateRelease(prInfos *models.StandardPrInfos, draft bool) (createdRelease *github.RepositoryRelease, err error) {
	relName := "Release " + prInfos.NextVersion
	defaultBranch := tools.GetDefaultBranch()
	fmt.Println(defaultBranch)
	releaseObject := github.RepositoryRelease{
		TargetCommitish: &defaultBranch,
		TagName:         &prInfos.NextVersion,
		Name:            &relName,
		Draft:           &draft,
	}
	createdRelease, _, err = GithubClient.Repositories.CreateRelease(
		ctx,
		prInfos.Owner,
		prInfos.Repo,
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
func PublishRelease(prInfos *models.StandardPrInfos, uploadArtifacts *string) (err error) {

	releaseIdStr, releaseIdBool := os.LookupEnv("ACI_RELEASE_ID")
	releaseId, err := strconv.ParseInt(releaseIdStr, 10, 64)
	if err != nil {
		fmt.Printf("%s of type %T", releaseIdStr, releaseIdStr)
		os.Exit(2)
	}

	draftFalse := false
	if !releaseIdBool {
		CreateRelease(prInfos, draftFalse)
	} else {
		existingRelease, _, err := GithubClient.Repositories.GetRelease(
			ctx,
			prInfos.Owner,
			prInfos.Repo,
			releaseId)
		if err != nil {
			return err
		}

		*existingRelease.Draft = draftFalse
		_, _, err = GithubClient.Repositories.EditRelease(
			ctx,
			prInfos.Owner,
			prInfos.Repo,
			releaseId,
			existingRelease)
		if err != nil {
			return err
		}
	}

	if uploadArtifacts != nil {
		filesAndInfos, err := tools.GetFilesAndInfos(uploadArtifacts)
		if err != nil {
			return err
		}

		for _, fileAndInfo := range filesAndInfos {
			log.Println("uploading file as asset to release", fileAndInfo)
			// Upload assets to GitHub Release
			_, _, err := GithubClient.Repositories.UploadReleaseAsset(
				ctx,
				prInfos.Owner,
				prInfos.Repo,
				releaseId,
				&github.UploadOptions{
					Name: fileAndInfo.Name(),
				},
				&fileAndInfo)
			if err != nil {
				log.Println("error at uploading asset to release: ", err)
			}
		}
	}
	return
}

// GetLatestReleaseVersion
func GetLatestReleaseVersion(owner string, repo string) (latestRelease *github.RepositoryRelease, err error) {
	latestRelease, _, err = GithubClient.Repositories.GetLatestRelease(ctx, owner, repo)
	return
}
