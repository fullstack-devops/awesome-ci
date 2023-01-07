package github

import (
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/semver"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
)

// GetPrInfos need the PullRequest-Number
func (ghrc *GitHubRichClient) GetPrInfos(prNumber int, mergeCommitSha string) (standardPrInfos *models.StandardPrInfos, prInfos *github.PullRequest, err error) {
	if prNumber != 0 {
		prInfos, _, err = ghrc.Client.PullRequests.Get(ctx, ghrc.Owner, ghrc.Repository, prNumber)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
	}
	if mergeCommitSha != "" && prNumber == 0 {
		prOpts := github.PullRequestListOptions{
			State:     "all",
			Sort:      "updated",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 10,
			},
		}
		pullRequests, _, err := ghrc.Client.PullRequests.List(ctx, ghrc.Owner, ghrc.Repository, &prOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
		var found int = 0
		for _, pr := range pullRequests {
			if pr.GetMergeCommitSHA() == mergeCommitSha {
				prInfos = pr
				found = found + 1
			}
		}
		if found > 1 {
			return nil, nil, fmt.Errorf("found more than one pull request, this should not be possible. please open an issue with all log files")
		}
	}

	if prInfos == nil {
		return nil, nil, fmt.Errorf("no pull request found, please check if all resources are specified")
	}

	// comment if awesome ci is running in CI mode
	isCI, isCIBool := os.LookupEnv("CI")
	_, isSilentBool := os.LookupEnv("ACI_SILENT")
	if isCIBool && !isSilentBool {
		if *prInfos.State == "open" && isCI == "true" {
			err = ghrc.CommentHelpToPullRequest(*prInfos.Number)
			if err != nil {
				log.Println(err)
			}
		}
	}

	prSHA := *prInfos.Head.SHA
	branchName := *prInfos.Head.Ref
	patchLevel := semver.ParsePatchLevel(branchName)

	var nextVersion = ""
	var latestVersion = ""
	// if an comment exists with aci_patch_level=major, make a major version!
	issueComments, err := ghrc.GetIssueComments(prNumber)
	if err != nil {
		return nil, nil, err
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
			aciVersionOverride := regexp.MustCompile(`^aci_version_override: ([0-9]+\.[0-9]+\.[0-9]+)`)
			aciPatchLevel := regexp.MustCompile(`^aci_patch_level: ([a-zA-Z]+)`)

			if aciVersionOverride.MatchString(*comment.Body) {
				nextVersion = aciVersionOverride.FindStringSubmatch(*comment.Body)[1]
				break
			}

			if aciPatchLevel.MatchString(*comment.Body) {
				patchLevel = semver.ParsePatchLevel(aciPatchLevel.FindStringSubmatch(*comment.Body)[1])
				break
			}

		}
	}

	if nextVersion == "" {
		repositoryRelease, err := ghrc.GetLatestReleaseVersion()
		if err == nil {
			latestVersion = *repositoryRelease.TagName
			nextVersion, err = semver.IncreaseVersion(patchLevel, latestVersion)
		} else {
			log.Traceln("no github release found")
			nextVersion, err = semver.IncreaseVersion(patchLevel, "0.0.0")
		}

		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Traceln("version override via pr comments specified")
	}

	standardPrInfos = &models.StandardPrInfos{
		PrNumber:       prNumber,
		Owner:          ghrc.Owner,
		Repo:           ghrc.Repository,
		BranchName:     branchName,
		Sha:            prSHA,
		ShaShort:       prSHA[:8],
		PatchLevel:     patchLevel,
		LatestVersion:  latestVersion,
		NextVersion:    nextVersion,
		MergeCommitSha: *prInfos.MergeCommitSHA,
	}
	return
}
