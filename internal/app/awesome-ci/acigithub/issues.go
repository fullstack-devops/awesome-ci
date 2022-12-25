package acigithub

import (
	"awesome-ci/internal/app/awesome-ci/tools"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v44/github"
)

var (
	direction, sort = "asc", "created"
)

func GetIssueComments(issueNumber int) (issueComments []*github.IssueComment, err error) {

	//FIXME: currently the GitHub API ignores Direction and Sort
	// so, we are using the default settings for the query and sorting the response afterwards
	var commentOpts = &github.IssueListCommentsOptions{
		Direction:   &direction,
		Sort:        &sort,
		ListOptions: standardListOptions,
	}
	issueComments, _, err = GithubClient.Issues.ListComments(ctx, owner, repo, issueNumber, commentOpts)

	if err == nil {
		var issueCommentsFiltered []*github.IssueComment
		for _, issueComment := range issueComments {
			if !strings.HasPrefix(*issueComment.Body, `<details><summary>Possible awesome-ci commands for this Pull Request</summary>`) {
				issueCommentsFiltered = append(issueCommentsFiltered, issueComment)
			}
		}

		icCount := len(issueCommentsFiltered)
		issueCommentsSorted := make([]*github.IssueComment, icCount)

		for i, issueComment := range issueCommentsFiltered {
			j := icCount - i - 1
			issueCommentsSorted[j] = issueComment
		}

		return issueCommentsSorted, err
	}

	return nil, err
}

func CommentHelpToPullRequest(number int) (err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}

	var commentOpts = &github.IssueListCommentsOptions{
		Direction:   &direction,
		Sort:        &sort,
		ListOptions: standardListOptions,
	}
	comments, _, err := GithubClient.Issues.ListComments(ctx, owner, repo, number, commentOpts)
	if err != nil {
		return fmt.Errorf("unable to list comments: %x", err)
	}

	body := `<details><summary>Possible awesome-ci commands for this Pull Request</summary>` +
		`</br>aci_patch_level: major</br>aci_version_override: 2.1.0` +
		`</br></br>Need more help?</br>Have a look at <a href="https://github.com/fullstack-devops/awesome-ci" target="_blank">my repo</a>` +
		`</br>This message was created by awesome-ci and can be disabled by the env variable <code>ACI_SILENT=true</code></details>`

	var prComment *github.IssueComment
	for _, prc := range comments {
		if strings.HasPrefix(*prc.Body, `<details><summary>Possible awesome-ci commands for this Pull Request</summary>`) {
			prComment = prc
		}
	}

	if prComment == nil {
		var prComment = &github.IssueComment{
			Body: &body,
		}

		err = CommentPullRequest(number, prComment)
	} else {
		// edit command if newer version deployed and commands changed
		if *prComment.Body != body {
			editedComment := &github.IssueComment{
				Body: &body,
			}
			_, _, err = GithubClient.Issues.EditComment(ctx, owner, repo, *prComment.ID, editedComment)
		}
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
