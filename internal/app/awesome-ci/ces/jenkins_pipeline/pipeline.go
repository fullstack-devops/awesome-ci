package jenkinspipeline

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
)

func DetectJenkinsPipeline() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isjenkinsUrl := os.LookupEnv("JENKINS_URL")

	serverUrl, isServerUrl := os.LookupEnv("GIT_URL")
	repository, isRepository := os.LookupEnv("GITHUB_REPOSITORY")
	token, isToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isjenkinsUrl && isServerUrl && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerUrl:  serverUrl,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no jenkins pipeline detected")
	}
}
