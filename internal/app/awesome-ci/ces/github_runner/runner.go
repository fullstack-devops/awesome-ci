package githubrunner

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
)

func DetectGitHubActionsRunner() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	serverUrl, isServerUrl := os.LookupEnv("GITHUB_SERVER_URL")
	repository, isRepository := os.LookupEnv("GITHUB_REPOSITORY")
	token, isToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isServerUrl && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerUrl:  serverUrl,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no github actions runner detected")
	}
}
