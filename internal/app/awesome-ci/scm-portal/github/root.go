package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
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
func NewGitHubClient(serverUrl *string, repoUrl *string, token *string) (githubRichClient *GitHubRichClient, err error) {
	gitHubTs := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	githubTc := oauth2.NewClient(context.Background(), gitHubTs)

	var githubClient *github.Client

	if !strings.HasPrefix(*serverUrl, "https://github.com") {
		githubClient, err = github.NewEnterpriseClient(*serverUrl, *serverUrl, githubTc)
		if err != nil {
			return nil, fmt.Errorf("error at initializing github client: %v", err)
		}
	} else {
		githubClient = github.NewClient(githubTc)
	}

	if rateLimit, _, err := githubClient.RateLimits(context.Background()); err != nil {
		return nil, fmt.Errorf("connection to GitHub could not be etablished: %v", err)
	} else {
		log.Tracef("remaining rates %d", rateLimit.Core.Remaining)
	}

	owner, repository := tools.DevideOwnerAndRepo(*repoUrl)

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
