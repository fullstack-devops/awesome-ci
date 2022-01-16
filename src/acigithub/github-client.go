package acigithub

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

var (
	GithubClient                       *github.Client
	githubServerUrl, isgithubServerUrl = os.LookupEnv("GITHUB_ENTERPRISE_SERVER_URL")
	githubToken, isgithubToken         = os.LookupEnv("GITHUB_TOKEN")
)

func devideOwnerAndRepo(fullRepo string) (owner string, repo string) {
	owner = strings.ToLower(strings.Split(fullRepo, "/")[0])
	repo = strings.ToLower(strings.Split(fullRepo, "/")[1])
	return
}

// NewGitHubClient Creates a new GitHub Client
// Needs the Environment Variables: GITHUB_TOKEN
// Needs the optional Environment Variables: GITHUB_ENTERPRISE_SERVER_URL
func NewGitHubClient() (githubClient *github.Client, err error) {
	if !isgithubToken {
		log.Fatalln("pleas set the GITHUB_TOKEN as environment variable!")
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
	GithubClient = githubClient
	return
}
