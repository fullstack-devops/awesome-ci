package acigithub

import (
	"awesome-ci/src/tools"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

var (
	GithubClient                         *github.Client
	ctx                                  = context.Background()
	githubServerUrl, isgithubServerUrl   = os.LookupEnv("GITHUB_ENTERPRISE_SERVER_URL")
	githubRepository, isgithubRepository = os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken           = os.LookupEnv("GITHUB_TOKEN")
	owner, repo                          string
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

	if !strings.HasSuffix(githubServerUrl, "/") {
		githubServerUrl = githubServerUrl + "/"
	}
	uploadUrl := fmt.Sprintf("%sapi/uploads/", githubServerUrl)

	if isgithubServerUrl {
		githubClient, err = github.NewEnterpriseClient(githubServerUrl, uploadUrl, githubTc)
		if err != nil {
			return nil, fmt.Errorf("error at initializing github client: %v", err)
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
