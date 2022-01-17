package acigithub

import (
	"awesome-ci/src/tools"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/v39/github"
)

// CreateRelease
func CreateRelease(version string, draft bool) (createdRelease *github.RepositoryRelease, err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := tools.DevideOwnerAndRepo(githubRepository)

	relName := "Release " + version
	defaultBranch := tools.GetDefaultBranch()

	releaseObject := github.RepositoryRelease{
		TargetCommitish: &defaultBranch,
		TagName:         &version,
		Name:            &relName,
		Draft:           &draft,
	}
	createdRelease, _, err = GithubClient.Repositories.CreateRelease(
		ctx,
		owner,
		repo,
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
func PublishRelease(version string, releaseId int64, uploadArtifacts *string) (err error) {
	draftFalse := false
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := tools.DevideOwnerAndRepo(githubRepository)

	if releaseId == 0 {
		releaseIdStr, releaseIdBool := os.LookupEnv("ACI_RELEASE_ID")
		if !releaseIdBool {
			log.Fatalln("Not release found creating one...")
			CreateRelease(version, draftFalse)
		}
		releaseId, err = strconv.ParseInt(releaseIdStr, 10, 64)
		if err != nil {
			fmt.Printf("%s of type %T", releaseIdStr, releaseIdStr)
			os.Exit(2)
		}
	}

	existingRelease, _, err := GithubClient.Repositories.GetRelease(
		ctx,
		owner,
		repo,
		releaseId)
	if err != nil {
		return err
	}

	*existingRelease.Draft = draftFalse
	_, _, err = GithubClient.Repositories.EditRelease(
		ctx,
		owner,
		repo,
		releaseId,
		existingRelease)
	if err != nil {
		return err
	}

	if uploadArtifacts != nil {
		filesAndInfos, err := tools.GetFilesAndInfos(uploadArtifacts)
		if err != nil {
			return err
		}

		for _, fileAndInfo := range filesAndInfos {
			fmt.Printf("uploading %s as asset to release", fileAndInfo.Name)
			// Upload assets to GitHub Release
			_, _, err := GithubClient.Repositories.UploadReleaseAsset(
				ctx,
				owner,
				repo,
				releaseId,
				&github.UploadOptions{
					Name: fileAndInfo.Name,
				},
				&fileAndInfo.File)
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
