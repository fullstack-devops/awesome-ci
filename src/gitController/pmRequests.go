package gitController

import (
	"awesome-ci/src/semver"
	"errors"
	"fmt"
	"strings"
)

/*
	Contains all functions for getting an Pull or Merge Request
*/

// GetPrInfos
func (ciEnv CIEnvironment) GetPrInfos(prNumber int) (aciPrInfos AciPrInfos, err error) {
	aciPrInfos.PrNumber = prNumber
	if ciEnv.Clients.GithubClient != nil {
		if prNumber != 0 {
			aciPrInfos.GithubPrInfos, _, err = ciEnv.Clients.GithubClient.PullRequests.Get(gitHubCtx, *ciEnv.GitInfos.Owner, *ciEnv.GitInfos.Repo, prNumber)
			if err != nil {
				err = errors.New(fmt.Sprintln("could not load any information about the current pull request", err))
				return
			}
		}

		prSHA := *aciPrInfos.GithubPrInfos.Head.SHA
		aciPrInfos.BranchName = *aciPrInfos.GithubPrInfos.Head.Ref
		aciPrInfos.Sha = prSHA
		aciPrInfos.ShaShort = prSHA[:8]

		branchName := *aciPrInfos.GithubPrInfos.Head.Ref
		aciPrInfos.PatchLevel = branchName[:strings.Index(branchName, "/")]

		// if an comment exists with aci=major, make a major version!
		issueComments, err := ciEnv.GetIssueComments(prNumber)
		if err != nil {
			panic(err)
		}
		for _, comment := range issueComments {
			// Must have permission in the repo to create a major version
			// MANNEQUIN|NONE https://docs.github.com/en/graphql/reference/enums#commentauthorassociation
			if strings.Contains(*comment.Body, "aci=major") && !strings.Contains("MANNEQUIN|NONE", *comment.AuthorAssociation) {
				aciPrInfos.PatchLevel = "major"
				break
			}
		}

		repositoryRelease, err := ciEnv.GetLatestReleaseVersion()
		if err != nil {
			return aciPrInfos, err
		}
		aciPrInfos.LatestVersion = *repositoryRelease.TagName
		aciPrInfos.CurrentVersion = *repositoryRelease.TagName
		aciPrInfos.NextVersion, err = semver.IncreaseVersion(aciPrInfos.PatchLevel, aciPrInfos.LatestVersion)
		if err != nil {
			return aciPrInfos, err
		}
	}
	if ciEnv.Clients.GitlabClient != nil {
		if prNumber != 0 {
			aciPrInfos.GitlabPrInfos, _, err = ciEnv.Clients.GitlabClient.MergeRequests.GetMergeRequest(1, prNumber, nil, nil)
			if err != nil {
				err = errors.New(fmt.Sprintln("could not load any information about the current pull request", err))
				return
			}
		}
		branchName := aciPrInfos.GitlabPrInfos.Reference
		aciPrInfos.PatchLevel = branchName[:strings.Index(branchName, "/")]
	}
	return
}
