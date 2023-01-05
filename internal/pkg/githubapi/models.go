package githubapi

import "github.com/google/go-github/v49/github"

type GitHubRichClient struct {
	Client     *github.Client
	Owner      string
	Repository string
}
