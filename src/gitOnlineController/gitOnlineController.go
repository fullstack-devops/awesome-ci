package gitOnlineController

import (
	"awesome-ci/src/models"
	"context"
	"log"

	"github.com/google/go-github/v39/github"
	"github.com/xanzy/go-gitlab"
)

var (
	gitHubCtx     = context.Background()
	CiEnvironment models.CIEnvironment
)

/* // GetPrNumberForBranch
func GetPrNumberForBranch(branch string) int {
	switch CiEnvironment.GitType {
	case "github":
		return github_getPrNumberForBranch(branch)
	}
	return 0
} */

func GetIssueComments(issueNumber int) (issueComments []*github.IssueComment, err error) {
	switch CiEnvironment.GitType {
	case "github":
		opts := github.IssueListCommentsOptions{}
		issueComments, _, err = CiEnvironment.GithubClient.Issues.ListComments(context.Background(), CiEnvironment.GitInfos.Repo, "", issueNumber, &opts)
	}
	return
}

// GetPrInfos
func GetPrInfos(prNumber int) (pullRequest *github.PullRequest, err error) {
	switch CiEnvironment.GitType {
	case "github":
		pullRequest, _, err = CiEnvironment.GithubClient.PullRequests.Get(gitHubCtx, CiEnvironment.GitInfos.Owner, CiEnvironment.GitInfos.Repo, prNumber)
	}
	return
}

// GetLatestReleaseVersion
func GetLatestReleaseVersion() (latestRelease *github.RepositoryRelease, err error) {
	switch CiEnvironment.GitType {
	case "github":
		latestRelease, _, err = CiEnvironment.GithubClient.Repositories.GetLatestRelease(gitHubCtx, CiEnvironment.GitInfos.Owner, CiEnvironment.GitInfos.Repo)
	}
	return
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
