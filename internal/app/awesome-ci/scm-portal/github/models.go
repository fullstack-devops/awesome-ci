package github

import "github.com/google/go-github/v70/github"

type GitHubRichClient struct {
	Client     *github.Client
	Owner      string
	Repository string
}
