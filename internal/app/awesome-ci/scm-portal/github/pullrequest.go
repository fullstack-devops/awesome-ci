package github

import (
	"awesome-ci/internal/pkg/semver"
	"fmt"

	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
)

// GetPrInfos need the PullRequest-Number
func (ghrc *GitHubRichClient) GetPrInfos(prNumber int, mergeCommitSha string) (prInfos *github.PullRequest, err error) {
	if prNumber != 0 {
		prInfos, _, err = ghrc.Client.PullRequests.Get(ctx, ghrc.Owner, ghrc.Repository, prNumber)
		if err != nil {
			return nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
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
			return nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
		var found int = 0
		for _, pr := range pullRequests {
			if pr.GetMergeCommitSHA() == mergeCommitSha {
				prInfos = pr
				found = found + 1
			}
		}
		if found > 1 {
			return nil, fmt.Errorf("found more than one pull request, this should not be possible. please open an issue with all log files")
		}
	}

	if prInfos == nil {
		return nil, fmt.Errorf("no pull request found, please check if all resources are specified")
	}

	// prSHA := *prInfos.Head.SHA
	branchName := *prInfos.Head.Ref
	patchLevel := semver.ParsePatchLevel(branchName)

	var nextVersion = ""
	var latestVersion = ""

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
			return nil, err
		}
	} else {
		log.Traceln("version override via pr comments specified")
	}

	return
}
