package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	"github.com/google/go-github/v70/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	ctx                 = context.Background()
	standardListOptions = github.ListOptions{
		PerPage: 100,
		Page:    1,
	}
)

// NewGitHubClient creates a new GitHub client
// Needs the ConnectCredentials
func NewGitHubClient(serverURL *string, repoURL *string, token *string) (githubRichClient *GitHubRichClient, err error) {
	gitHubTS := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	githubTc := oauth2.NewClient(context.Background(), gitHubTS)

	var githubClient *github.Client

	if !strings.HasPrefix(*serverURL, "https://github.com") {
		githubClient, err = github.NewClient(githubTc).WithEnterpriseURLs(*serverURL, *serverURL)
		if err != nil {
			return nil, fmt.Errorf("error at initializing github client: %v", err)
		}
	} else {
		githubClient = github.NewClient(githubTc)
	}

	if rateLimit, _, err := githubClient.RateLimit.Get(context.Background()); err != nil {
		if !strings.Contains(fmt.Sprintf("%v", err), "404 Rate limiting is not enabled") {
			logrus.Traceln("rate limiting is not enabled, but connection is good")
			return nil, fmt.Errorf("connection to GitHub could not be etablished: %v", err)
		}
	} else {
		logrus.Tracef("remaining rates %d", rateLimit.Core.Remaining)
	}

	owner, repository := tools.DevideOwnerAndRepo(*repoURL)

	return &GitHubRichClient{
		Client:     githubClient,
		Owner:      owner,
		Repository: repository,
	}, nil
}
