package githubrunner

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
)

func DetectGitHubActionsRunner() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	serverURL, isServerURL := os.LookupEnv("GITHUB_SERVER_URL")
	repository, isRepository := os.LookupEnv("GITHUB_REPOSITORY")
	token, isToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isServerURL && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerURL:  serverURL,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no github actions runner detected")
	}
}
