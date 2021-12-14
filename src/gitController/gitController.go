package gitController

import (
	"context"

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
		issueComments, _, err = ciEnv.Clients.GithubClient.Issues.ListComments(gitHubCtx, *ciEnv.GitInfos.Owner, *ciEnv.GitInfos.Repo, issueNumber, &opts)
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
