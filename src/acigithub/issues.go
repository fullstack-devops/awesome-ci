package acigithub

import (
	"awesome-ci/src/tools"
	"fmt"
	"log"

	"github.com/google/go-github/v39/github"
)

var (
	direction, sort = "asc", "created"
)

func GetIssueComments(issueNumber int, owner string, repo string) (issueComments []*github.IssueComment, err error) {
	var commentOpts = &github.IssueListCommentsOptions{
		Direction: &direction,
		// Sort:      &sort,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	issueComments, _, err = GithubClient.Issues.ListComments(ctx, owner, repo, issueNumber, commentOpts)
	return
}

func CommentHelpToPullRequest(number int) (err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := tools.DevideOwnerAndRepo(githubRepository)

	var commentOpts = &github.IssueListCommentsOptions{
		Direction: &direction,
		Sort:      &sort,
		ListOptions: github.ListOptions{
			PerPage: 30,
			Page:    1,
		},
	}
	comments, _, err := GithubClient.Issues.ListComments(ctx, owner, repo, number, commentOpts)
	if err != nil {
		return fmt.Errorf("unable to list comments: %x", err)
	}

	body := `<details><summary>Possible awesome-ci commands for this Pull Request</summary>` +
		`</br>aci_patch_level: major</br>aci_version_override: 2.1.0` +
		`</br></br>Need more help?</br>Have a look at <a href="https://github.com/fullstack-devops/awesome-ci" target="_blank">my repo</a></details>`

	var myCommand *github.IssueComment
	for _, prc := range comments {
		if body == *prc.Body {
			myCommand = prc
		}
	}

	if myCommand == nil {
		var prComment = &github.IssueComment{
			Body: &body,
		}

		err = CommentPullRequest(number, prComment)
	}
	return
}

func CommentPullRequest(number int, comment *github.IssueComment) (err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := tools.DevideOwnerAndRepo(githubRepository)

	_, _, err = GithubClient.Issues.CreateComment(ctx, owner, repo, number, comment)
	return
}
