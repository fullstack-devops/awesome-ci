package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"
	"github.com/sirupsen/logrus"

	"github.com/google/go-github/v64/github"
)

var (
	direction, sort = "asc", "created"
	detailsOpening  = `<details><summary>Possible awesome-ci commands for this Pull Request</summary>`

	// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	semVerRegex = `(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
)

func (ghrc *GitHubRichClient) GetIssueComments(issueNumber int) (issueComments []*github.IssueComment, err error) {

	//FIXME: currently the GitHub API ignores Direction and Sort
	// so, we are using the default settings for the query and sorting the response afterwards
	var commentOpts = &github.IssueListCommentsOptions{
		Direction:   &direction,
		Sort:        &sort,
		ListOptions: standardListOptions,
	}
	issueComments, _, err = ghrc.Client.Issues.ListComments(ctx, ghrc.Owner, ghrc.Repository, issueNumber, commentOpts)

	if err == nil {
		var issueCommentsFiltered []*github.IssueComment
		for _, issueComment := range issueComments {
			if !strings.HasPrefix(*issueComment.Body, detailsOpening) {
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

func (ghrc *GitHubRichClient) CommentHelpToPullRequest(number int) (err error) {
	var commentOpts = &github.IssueListCommentsOptions{
		Direction:   &direction,
		Sort:        &sort,
		ListOptions: standardListOptions,
	}
	comments, _, err := ghrc.Client.Issues.ListComments(ctx, ghrc.Owner, ghrc.Repository, number, commentOpts)
	if err != nil {
		return fmt.Errorf("unable to list comments: %x", err)
	}

	body := detailsOpening + `
	| command | description |
	| --- | --- |
	| aci_patch_level: major | create a major version bump |
	| aci_version_override: 2.1.0 | set the version to 2.1.0 using only semver compatible syntax! |

	Need more help?

	Have a look at [my repo](https://github.com/fullstack-devops/awesome-ci)

	This message was created by awesome-ci and can be disabled by the env variable <code>ACI_SILENT=true</code>
	</details>`

	var prComment *github.IssueComment
	for _, prc := range comments {
		if strings.HasPrefix(*prc.Body, detailsOpening) {
			prComment = prc
		}
	}

	if prComment == nil {
		var prComment = &github.IssueComment{
			Body: &body,
		}

		_, _, err = ghrc.Client.Issues.CreateComment(ctx, ghrc.Owner, ghrc.Repository, number, prComment)
		// err = CommentPullRequest(number, prComment)
	} else {
		// edit command if newer version deployed and commands changed
		if *prComment.Body != body {
			editedComment := &github.IssueComment{
				Body: &body,
			}
			_, _, err = ghrc.Client.Issues.EditComment(ctx, ghrc.Owner, ghrc.Repository, *prComment.ID, editedComment)
		}
	}
	return
}

func (ghrc *GitHubRichClient) SearchIssuesForOverrides(prNumber int) (nextVersion *string, patchLevel *semver.PatchLevel, err error) {
	var pLevel semver.PatchLevel
	// if an comment exists with aci_patch_level=major, make a major version!
	issueComments, err := ghrc.GetIssueComments(prNumber)
	if err != nil {
		return
	}

	for _, comment := range issueComments {
		// FIXME: access must be restricted but GITHUB_TOKEN doesn't get informations.
		// Refs: https://docs.github.com/en/rest/collaborators/collaborators#list-repository-collaborators
		// Refs: https://docs.github.com/en/actions/security-guides/automatic-token-authentication#permissions-for-the-github_token

		// Must be a collaborator to have permission to create an override
		// isCollaborator, resp, err := GithubClient.Repositories.IsCollaborator(ctx, owner, repo, *comment.User.Login)
		// if err != nil {
		// 	return nil, nil, err
		// }
		// fmt.Println(resp.StatusCode)

		// if isCollaborator {
		if true {
			aciVersionOverride := regexp.MustCompile(`^aci_version_override: (` + semVerRegex + `)`)
			aciPatchLevel := regexp.MustCompile(`^aci_patch_level: ([a-zA-Z]+)`)

			if aciVersionOverride.MatchString(*comment.Body) {
				nextVersion = &aciVersionOverride.FindStringSubmatch(*comment.Body)[1]
				logrus.Infof("found version override in comment form %s at %s", comment.User.GetLogin(), comment.GetCreatedAt())
				break
			}

			if aciPatchLevel.MatchString(*comment.Body) {
				pLevel, err = semver.ParsePatchLevel(aciPatchLevel.FindStringSubmatch(*comment.Body)[1])
				if err != nil && err != semver.ErrUseMinimalPatchVersion {
					logrus.Warnln(err)
				}
				logrus.Infof("found patchLevel override in comment form %s at %s", comment.User.GetLogin(), comment.GetCreatedAt())
				patchLevel = &pLevel
				break
			}
		}
	}
	return
}

/* func (ghrc *GitHubRichClient) CommentPullRequest(number int, comment *github.IssueComment) (err error) {
	_, _, err = ghrc.Client.Issues.CreateComment(ctx, ghrc.Owner, ghrc.Repository, number, comment)
	return
} */
