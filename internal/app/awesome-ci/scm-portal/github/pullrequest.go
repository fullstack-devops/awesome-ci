package github

import (
	"fmt"

	"github.com/google/go-github/v67/github"
	"github.com/sirupsen/logrus"
)

// GetPrInfos need the PullRequest-Number
func (ghrc *GitHubRichClient) GetPrInfos(prNumber int, mergeCommitSha string) (prInfos *github.PullRequest, err error) {
	if prNumber != 0 {

		prInfos, _, err = ghrc.Client.PullRequests.Get(ctx, ghrc.Owner, ghrc.Repository, prNumber)
		if err != nil {
			return nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
		logrus.Infof("found pull request '%s' with given number %d", *prInfos.Title, prNumber)

		return
	} else if mergeCommitSha != "" {

		logrus.Infoln("listing pull requests to compare with merge commit sha")
		prOpts := github.PullRequestListOptions{
			State:     "all",
			Sort:      "updated",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 50,
			},
		}
		pullRequests, _, err := ghrc.Client.PullRequests.List(ctx, ghrc.Owner, ghrc.Repository, &prOpts)
		if err != nil {
			return nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}

		logrus.Infof("listed %d pull requests, searching for maching mergeCommitSha", len(pullRequests))
		var found = 0
		for _, pr := range pullRequests {
			if pr.GetMergeCommitSHA() == mergeCommitSha {
				logrus.Infof("found matching pull requests with number %d, and mergeCommitSha %s", *pr.Number, *pr.MergeCommitSHA)
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
	return
}
