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

// GetLatestReleaseVersion
func GetLatestReleaseVersion() string {
	switch CiEnvironment.GitType {
	case "github":
		return github_getLatestReleaseVersion()
	}
	return ""
}

// CreateNextGitHubRelease
func CreateNextGitHubRelease(releaseBranch string, newReleaseVersion string, preRelease *bool, uploadArtifacts string) {
	switch CiEnvironment.GitType {
	case "github":
		github_createNextGitHubRelease(
			releaseBranch,
			newReleaseVersion,
			*preRelease,
			uploadArtifacts)
	}
}
