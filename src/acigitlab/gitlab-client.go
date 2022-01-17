package acigitlab

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

var (
	GitLabClient                   *gitlab.Client
	gitlabCiToken, isgitlabCiToken = os.LookupEnv("CI_JOB_TOKEN")
)

// NewGitHubClient Creates a new GitHub Client
// Needs the Environment Variables: GITHUB_TOKEN
// Needs the optional Environment Variables: GITHUB_ENTERPRISE_SERVER_URL
func NewGitHubClient() (gitlabClient *gitlab.Client, err error) {
	if isgitlabCiToken {
		gitlabClient, err = gitlab.NewClient(gitlabCiToken)
		if err != nil {
			fmt.Errorf("error at initializing gitlab client: %v", err)
		}
	} else {
		panic("Not running in GitLab CI?!")
	}
	GitLabClient = gitlabClient
	return
}
