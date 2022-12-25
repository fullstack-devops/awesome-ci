package githubapi

import (
	"awesome-ci/internal/pkg/tools"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

var (
	GithubClient                         *github.Client
	ctx                                  = context.Background()
	githubServerUrl, isgithubServerUrl   = os.LookupEnv("GITHUB_SERVER_URL")
	githubRepository, isgithubRepository = os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken           = os.LookupEnv("GITHUB_TOKEN")
	owner, repo                          string
	standardListOptions                  = github.ListOptions{
		PerPage: 100,
		Page:    1,
	}
)

// NewGitHubClient Creates a new GitHub Client
// Needs the Environment Variables: GITHUB_TOKEN
// Needs the optional Environment Variables: GITHUB_ENTERPRISE_SERVER_URL
func NewGitHubClient() (githubClient *github.Client, err error) {
	if !isgithubToken {
		log.Fatalln("please set the GITHUB_TOKEN as environment variable!")
	}

	gitHubTs := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	githubTc := oauth2.NewClient(context.Background(), gitHubTs)

	if isgithubServerUrl {
		if githubServerUrl != "https://github.com" {
			githubClient, err = github.NewEnterpriseClient(githubServerUrl, githubServerUrl, githubTc)
			if err != nil {
				return nil, fmt.Errorf("error at initializing github client: %v", err)
			}
		} else {
			githubClient = github.NewClient(githubTc)
		}
	} else {
		githubClient = github.NewClient(githubTc)
	}

	if isgithubRepository {
		owner, repo = tools.DevideOwnerAndRepo(githubRepository)
	}

	GithubClient = githubClient
	return
}
