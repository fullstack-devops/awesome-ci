package githubapi

import "github.com/google/go-github/v44/github"

type GitHubRichClient struct {
	Client     *github.Client
	Owner      string
	Repository string
}

type ConnectCredentials struct {
	ServerUrl  string
	Repository string
	Token      string
}
