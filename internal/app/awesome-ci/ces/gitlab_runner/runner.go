package gitlabrunner

import (
	"awesome-ci/internal/pkg/models"
	"fmt"
	"os"
)

func DetectGitLabActionsRunner() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	serverUrl, isServerUrl := os.LookupEnv("CI_SERVER_URL")
	repository, isRepository := os.LookupEnv("CI_PROJECT_URL")
	token, isToken := os.LookupEnv("CI_JOB_TOKEN")

	if isCI && isServerUrl && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerUrl:  serverUrl,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no gitlab runner detected")
	}
}
