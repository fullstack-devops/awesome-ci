package jenkinspipeline

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
)

func DetectJenkinsPipeline() (connects *models.ConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isjenkinsURL := os.LookupEnv("JENKINS_URL")

	serverURL, isServerURL := os.LookupEnv("GIT_URL")
	repository, isRepository := os.LookupEnv("GITHUB_REPOSITORY")
	token, isToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isjenkinsURL && isServerURL && isRepository && isToken {
		return &models.ConnectCredentials{
			ServerURL:  serverURL,
			Repository: repository,
			Token:      token,
		}, nil

	} else {
		return nil, fmt.Errorf("no jenkins pipeline detected")
	}
}
