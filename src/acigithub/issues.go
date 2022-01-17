package acigithub

import (
	"github.com/google/go-github/v39/github"
)

func GetIssueComments(issueNumber int, owner string, repo string) (issueComments []*github.IssueComment, err error) {
	// opts := github.IssueListCommentsOptions{}
	issueComments, _, err = GithubClient.Issues.ListComments(ctx, owner, repo, issueNumber, nil)
	return
}
