package detect

import (
	"awesome-ci/internal/pkg/models"
	"fmt"
	"os"
)

func CheckGitHubActionsRunner() (isRunner bool, connects models.StandardConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent
	githubServerUrl, isgithubServerUrl := os.LookupEnv("GITHUB_SERVER_URL")
	githubRepository, isgithubRepository := os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isgithubServerUrl && isgithubRepository && isgithubToken {
		return true, models.StandardConnectCredentials{
			ServerUrl:  githubServerUrl,
			Repository: githubRepository,
			Token:      githubToken,
		}, nil

	} else {
		return false, models.StandardConnectCredentials{}, fmt.Errorf("no github actions runner detected")
	}
}
