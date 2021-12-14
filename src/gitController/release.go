package gitController

import (
	"awesome-ci/src/controlEnvs"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/v39/github"
)

// CreateRelease
func (ciEnv CIEnvironment) CreateRelease(aciPrInfos *AciPrInfos, draft bool) (createdRelease *github.RepositoryRelease, err error) {
	relName := "Release " + aciPrInfos.NextVersion

	if ciEnv.Clients.GithubClient != nil {
		releaseObject := github.RepositoryRelease{
			TargetCommitish: &aciPrInfos.BranchName,
			TagName:         &aciPrInfos.NextVersion,
			Name:            &relName,
			Draft:           &draft,
		}
		createdRelease, _, err = ciEnv.Clients.GithubClient.Repositories.CreateRelease(
			gitHubCtx,
			*ciEnv.GitInfos.Owner,
			*ciEnv.GitInfos.Repo,
			&releaseObject)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(createdRelease)

		envVariables, err := controlEnvs.OpenEnvFile(ciEnv.RunnerInfo.EnvFile)
		if err != nil {
			return nil, err
		}
		envVariables.Set("ACI_RELEASE_ID", fmt.Sprintf("d", *createdRelease.ID))
		err = envVariables.SaveEnvFile()
	}
	if ciEnv.Clients.GitlabClient != nil {
		/* releaseObject := gitlab.Release{
			TagName: aciPrInfos.NextVersion,
			Name:    relName,
		} */
		log.Println("Creating a release to GitLab is not jet implemented")
	}
	return
}

// PublishRelease
func (ciEnv CIEnvironment) PublishRelease(aciPrInfos *AciPrInfos, uploadArtifacts *string) (err error) {

	releaseIdStr, releaseIdBool := os.LookupEnv("ACI_RELEASE_ID")
	releaseId, err := strconv.ParseInt(releaseIdStr, 10, 64)
	if err != nil {
		fmt.Printf("%s of type %T", releaseIdStr, releaseIdStr)
		os.Exit(2)
	}

	if ciEnv.Clients.GithubClient != nil {
		draftFalse := false
		if !releaseIdBool {
			ciEnv.CreateRelease(aciPrInfos, draftFalse)
		} else {
			existingRelease, _, err := ciEnv.Clients.GithubClient.Repositories.GetRelease(
				gitHubCtx,
				*ciEnv.GitInfos.Owner,
				*ciEnv.GitInfos.Repo,
				releaseId)
			if err != nil {
				return err
			}

			existingRelease.Draft = &draftFalse
			_, _, err = ciEnv.Clients.GithubClient.Repositories.EditRelease(
				gitHubCtx,
				*ciEnv.GitInfos.Owner,
				*ciEnv.GitInfos.Repo,
				releaseId,
				existingRelease)
			if err != nil {
				return err
			}
		}

		if uploadArtifacts != nil {
			filesAndInfos, err := getFilesAndInfos(uploadArtifacts)
			if err != nil {
				return err
			}

			for _, fileAndInfo := range filesAndInfos {
				log.Println("uploading file as asset to release", fileAndInfo)
				// Upload assets to GitHub Release
				_, _, err := ciEnv.Clients.GithubClient.Repositories.UploadReleaseAsset(
					gitHubCtx,
					*ciEnv.GitInfos.Owner,
					*ciEnv.GitInfos.Repo,
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
	}
	if ciEnv.Clients.GitlabClient != nil {
		/* releaseObject := gitlab.Release{
			TagName: aciPrInfos.NextVersion,
			Name:    relName,
		} */
		log.Println("Publishing a release to GitLab is not jet implemented")
	}
	return
}
