package gitlabrunner

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
)

func DetectGitLabActionsRunner() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	serverURL, isServerURL := os.LookupEnv("CI_SERVER_URL")
	repository, isRepository := os.LookupEnv("CI_PROJECT_URL")
	token, isToken := os.LookupEnv("CI_JOB_TOKEN")

	if isCI && isServerURL && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerURL:  serverURL,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no gitlab runner detected")
	}
}
