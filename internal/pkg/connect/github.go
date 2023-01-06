package connect

import (
	"awesome-ci/internal/pkg/detect"
	"awesome-ci/internal/pkg/tools"
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// NewGitHubClient creates a new GitHub client
// Needs the ConnectCredentials
func NewGitHubClient(creds ConnectCredentials) (githubRichClient *GitHubRichClient, err error) {
	gitHubTs := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: creds.TokenPlain},
	)
	githubTc := oauth2.NewClient(context.Background(), gitHubTs)

	var githubClient *github.Client

	if !strings.HasPrefix(creds.ServerUrl, "https://github.com") {
		githubClient, err = github.NewEnterpriseClient(creds.ServerUrl, creds.ServerUrl, githubTc)
		if err != nil {
			return nil, fmt.Errorf("error at initializing github client: %v", err)
		}
	} else {
		githubClient = github.NewClient(githubTc)
	}

	if rateLimit, _, err := githubClient.RateLimits(context.Background()); err != nil {
		return nil, fmt.Errorf("connection to GitHub could not be etablished: %v", err)
	} else {
		log.Tracef("remaining rates %", rateLimit.Core.Remaining)
	}

	owner, repository := tools.DevideOwnerAndRepo(creds.Repository)

	return &GitHubRichClient{
		Client:     githubClient,
		Owner:      owner,
		Repository: repository,
	}, nil
}

func (ghrc *GitHubRichClient) TestGitHubClientConnection() (err error) {
	_, _, err = ghrc.Client.RateLimits(context.Background())
	return
}

func ConnectToGitHub() (githubRichClient *GitHubRichClient, err error) {
	if tools.CheckFileExists(rcFileName) {
		log.Infof("Connecting to Github via local %s file", rcFileName)
		return ConnectWithRcFileToGithub()
	} else {
		isGithubActionsRunner, _, _ := detect.CheckGitHubActionsRunner()
		if isGithubActionsRunner {
			log.Info("Connecting to Github via GitHub actions runner")
			return ConnectWithActionsToGitHub()
		}
		isJenkinsPipeline, _, _ := detect.CheckJenkinsPipeline()
		if isJenkinsPipeline {
			log.Info("Connecting to Github via Jenkins pipeline")
			return ConnectWithJenkinsToGitHub()
		}
	}
	return nil, fmt.Errorf("could not connect with any method with GitHub. Please read the Docs")
}

func UpdateRcFileForGithub(server string, repo string, token string) {
	rcFile := NewRcInstance()

	rcFile.UpdateServerType(serverTypeGitHub)
	if err := rcFile.UpdateCreds(server, repo, token); err != nil {
		log.Fatalln(err)
	}

	// testing connection
	if _, err := NewGitHubClient(rcFile.ConnectCreds); err != nil {
		log.Fatalln(err)
	}

	rcFile.Save()
}

func ConnectWithRcFileToGithub() (githubRichClient *GitHubRichClient, err error) {
	rcFile := NewRcInstance()

	if rcFile.Exists() {
		creds, err := rcFile.Load()
		if err != nil {
			log.Fatalln(err)
		}

		return NewGitHubClient(*creds)

	} else {

		return nil, fmt.Errorf("no %s file found", rcFileName)
	}
}

// GitHub Actions
func ConnectWithActionsToGitHub() (githubRichClient *GitHubRichClient, err error) {
	isRunner, creds, _ := detect.CheckGitHubActionsRunner()
	if isRunner {
		return NewGitHubClient(creds)
	} else {
		err = fmt.Errorf("")
		return
	}
}

// GitHub Actions
func ConnectWithJenkinsToGitHub() (githubRichClient *GitHubRichClient, err error) {
	isRunner, creds, _ := detect.CheckJenkinsPipeline()
	if isRunner {
		return NewGitHubClient(creds)
	} else {
		err = fmt.Errorf("")
		return
	}
}
