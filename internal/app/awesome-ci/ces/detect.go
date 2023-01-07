package ces

import (
	githubrunner "awesome-ci/internal/app/awesome-ci/ces/github_runner"
	gitlabrunner "awesome-ci/internal/app/awesome-ci/ces/gitlab_runner"
	jenkinspipeline "awesome-ci/internal/app/awesome-ci/ces/jenkins_pipeline"
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/rcpersist"
	"os"
)

// DetectCes detects the current "code execution side"
func DetectCes() (cesType rcpersist.CESType,
	scmPortalType rcpersist.SCMPortalType,
	connCreds *models.ConnectCredentials,
	err error) {

	// try github runner
	if connCreds, _ = githubrunner.DetectGitHubActionsRunner(); connCreds != nil {
		cesType = rcpersist.CESTypeGitHubRunner
		scmPortalType = rcpersist.SCMPortalTypeGitHub
		return
	}

	// try gitlab runner
	if connCreds, _ = gitlabrunner.DetectGitLabActionsRunner(); connCreds != nil {
		cesType = rcpersist.CESTypeGitLabRunner
		scmPortalType = rcpersist.SCMPortalTypeGitLab
		return
	}

	// try jenkins pipeline
	if connCreds, _ = jenkinspipeline.DetectJenkinsPipeline(); connCreds != nil {
		cesType = rcpersist.CESTypeJenkinsPipeline
		scmPortalType = rcpersist.SCMPortalTypeGitHub
		return
	}

	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "" && ciPresent

	if !isCI {
		cesType = rcpersist.CESTypeLoMa
		return
	}

	return
}
