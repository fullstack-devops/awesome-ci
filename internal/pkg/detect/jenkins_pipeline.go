package detect

import (
	"awesome-ci/internal/pkg/models"
	"fmt"
	"os"
)

func CheckJenkinsPipeline() (isJenkins bool, connects models.StandardConnectCredentials, err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent
	_, isjenkinsUrl := os.LookupEnv("JENKINS_URL")
	githubServerUrl, isgithubServerUrl := os.LookupEnv("GIT_URL")
	githubRepository, isgithubRepository := os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken := os.LookupEnv("GITHUB_TOKEN")

	if isCI && isjenkinsUrl && isgithubServerUrl && isgithubRepository && isgithubToken {
		return true, models.StandardConnectCredentials{
			ServerUrl:  githubServerUrl,
			Repository: githubRepository,
			Token:      githubToken,
		}, nil

	} else {
		return false, models.StandardConnectCredentials{}, fmt.Errorf("no jenkins pipeline detected")
	}
}
