package gitController

import (
	"context"
	"log"

	"github.com/google/go-github/v39/github"
	"github.com/xanzy/go-gitlab"
)

var (
	gitHubCtx = context.Background()
)

type CIEnvironment struct {
	GitInfos struct {
		ApiUrl            *string
		ApiToken          *string
		Repo              *string
		Owner             *string
		IsOrg             *bool
		DefaultBranchName *string
	}
	Clients struct {
		GithubClient *github.Client
		GitlabClient *gitlab.Client
	}
	RunnerType string
	RunnerInfo struct {
		EnvFile string
	}
}

// GetIssueComments
func (ciEnv CIEnvironment) GetIssueComments(issueNumber int) (issueComments []*github.IssueComment, err error) {
	if ciEnv.Clients.GithubClient != nil {
		opts := github.IssueListCommentsOptions{}
		issueComments, _, err = ciEnv.Clients.GithubClient.Issues.ListComments(gitHubCtx, *ciEnv.GitInfos.Repo, "", issueNumber, &opts)
	}
	return
}

// GetPrInfos
func (ciEnv CIEnvironment) GetPrInfos(prNumber int) (pullRequest *github.PullRequest, err error) {
	if ciEnv.Clients.GithubClient != nil {
		pullRequest, _, err = ciEnv.Clients.GithubClient.PullRequests.Get(gitHubCtx, *ciEnv.GitInfos.Owner, *ciEnv.GitInfos.Repo, prNumber)
	}
	return
}

// GetLatestReleaseVersion
func (ciEnv CIEnvironment) GetLatestReleaseVersion() (latestRelease *github.RepositoryRelease, err error) {
	if ciEnv.Clients.GithubClient != nil {
		latestRelease, _, err = ciEnv.Clients.GithubClient.Repositories.GetLatestRelease(gitHubCtx, *ciEnv.GitInfos.Owner, *ciEnv.GitInfos.Repo)
	}
	return
}

// CreateNextGitHubRelease
func (ciEnv CIEnvironment) ManageGitRelease(releaseObject interface{}, uploadArtifacts *string) (err error) {
	switch rel := releaseObject.(type) {
	case github.RepositoryRelease:
		repositoryRelease, _, err := ciEnv.Clients.GithubClient.Repositories.CreateRelease(
			gitHubCtx,
			*ciEnv.GitInfos.Owner,
			*ciEnv.GitInfos.Repo,
			&rel)
		if err != nil {
			return err
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
					repositoryRelease.GetID(),
					&github.UploadOptions{
						Name: fileAndInfo.Name(),
					},
					&fileAndInfo)
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
