package githubapi

import (
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/tools"
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

var (
	ctx                 = context.Background()
	standardListOptions = github.ListOptions{
		PerPage: 100,
		Page:    1,
	}
)

// NewGitHubClient Creates a new GitHub Client
// Needs the Environment Variables: GITHUB_TOKEN
// Needs the optional Environment Variables: GITHUB_ENTERPRISE_SERVER_URL
func NewGitHubClient(creds models.StandardConnectCredentials) (githubRichClient *GitHubRichClient, err error) {
	gitHubTs := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: creds.Token},
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

	owner, repository := tools.DevideOwnerAndRepo(creds.Repository)

	return &GitHubRichClient{
		Client:     githubClient,
		Owner:      owner,
		Repository: repository,
	}, nil
}
