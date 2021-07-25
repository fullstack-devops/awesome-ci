package gitOnlineController

import "awesome-ci/models"

var CiEnvironment models.CIEnvironment

// GetPrNumberForBranch
func GetPrNumberForBranch(branch string) int {
	switch CiEnvironment.GitType {
	case "github":
		return github_getPrNumberForBranch(branch)
	}
	return 0
}

func GetIssueComments(issueNumber int) (issueComments []models.GitHubIssueComment, err error) {
	switch CiEnvironment.GitType {
	case "github":
		issueComments, err = github_getIssueComments(issueNumber)
	}
	return
}

// GetPrNumberForBranch
func GetPrInfos(prNumber int) (prInfos models.GitHubPullRequest, err error) {
	switch CiEnvironment.GitType {
	case "github":
		prInfos, err = github_getPrInfos(prNumber)
	}
	return
}

// GetLatestReleaseVersion
func GetLatestReleaseVersion() string {
	switch CiEnvironment.GitType {
	case "github":
		return github_getLatestReleaseVersion()
	}
	return ""
}

// CreateNextGitHubRelease
func CreateNextGitHubRelease(releaseBranch string, newReleaseVersion string, preRelease *bool, isDryRun *bool, uploadArtifacts *string) {
	switch CiEnvironment.GitType {
	case "github":
		github_createNextGitHubRelease(
			releaseBranch,
			newReleaseVersion,
			preRelease,
			isDryRun,
			uploadArtifacts)
	}
}
