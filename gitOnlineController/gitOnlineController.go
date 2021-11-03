package gitOnlineController

import (
	"awesome-ci/models"
	"context"
	"log"

	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
)

var (
	gitHubCtx     = context.Background()
	CiEnvironment models.CIEnvironment
)

// GetPrNumberForBranch
func GetPrNumberForBranch(branch string) int {
	switch CiEnvironment.GitType {
	case "github":
		return github_getPrNumberForBranch(branch)
	}
	return 0
}

func GetIssueComments(issueNumber int) (issueComments []models.GitHubIssueComment, err error) {
	switch CiEnvironment.GitType {
	case "github":
		issueComments, err = github_getIssueComments(issueNumber)
	}
	return
}

// GetPrNumberForBranch
func GetPrInfos(prNumber int) (pullRequest *github.PullRequest, err error) {
	switch CiEnvironment.GitType {
	case "github":
		pullRequest, _, err = CiEnvironment.GithubClient.PullRequests.Get(gitHubCtx, CiEnvironment.GitInfos.Owner, CiEnvironment.GitInfos.Repo, prNumber)
	}
	return
}

// GetLatestReleaseVersion
func GetLatestReleaseVersion() string {
	switch CiEnvironment.GitType {
	case "github":
		return github_getLatestReleaseVersion()
	}
	return ""
}

// CreateNextGitHubRelease
func CreateNextGitRelease(releaseObject interface{}, uploadArtifacts *string) (err error) {
	switch rel := releaseObject.(type) {
	case github.RepositoryRelease:
		repositoryRelease, _, err := CiEnvironment.GithubClient.Repositories.CreateRelease(
			gitHubCtx,
			CiEnvironment.GitInfos.Owner,
			CiEnvironment.GitInfos.Repo,
			&rel)
		if err != nil {
			return err
		}

		if uploadArtifacts != nil {
			filesAndInfos, err := getFilesAndInfos(*uploadArtifacts)
			if err != nil {
				return err
			}

			for _, fileAndInfo := range filesAndInfos {
				log.Println("uploading file as asset to release", fileAndInfo)
				// Upload assets to GitHub Release
				_, _, err := CiEnvironment.GithubClient.Repositories.UploadReleaseAsset(
					gitHubCtx,
					CiEnvironment.GitInfos.Owner,
					CiEnvironment.GitInfos.Repo,
					repositoryRelease.GetID(),
					&github.UploadOptions{
						Name: fileAndInfo.Name,
					},
					&fileAndInfo.File)
				if err != nil {
					log.Println("error at uploading asset to release: ", err)
				}
			}
		}
	case *gitlab.Release:
		log.Println("Creating a release to GitLab is not jet implemented")
	}
	return
}
