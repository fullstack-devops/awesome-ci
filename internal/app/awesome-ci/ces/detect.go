package ces

import (
	githubrunner "awesome-ci/internal/app/awesome-ci/ces/github_runner"
	gitlabrunner "awesome-ci/internal/app/awesome-ci/ces/gitlab_runner"
	jenkinspipeline "awesome-ci/internal/app/awesome-ci/ces/jenkins_pipeline"
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/rcpersist"
	"os"

	log "github.com/sirupsen/logrus"
)

// DetectCes detects the current "code execution side"
// if DetectCes finds an .rc file, this will be used first
func DetectCes() (cesType rcpersist.CESType,
	scmPortalType rcpersist.SCMPortalType,
	connCreds *models.ConnectCredentials,
	err error) {

	rcFile := rcpersist.NewRcInstance()

	if creds, errRc := rcFile.Load(); errRc == nil {
		cesType = rcFile.CESType
		scmPortalType = rcFile.SCMPortalType
		connCreds = &models.ConnectCredentials{
			ServerUrl:  creds.ServerUrl,
			Repository: creds.Repository,
			Token:      *creds.TokenPlain,
		}
		return
	}

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
		log.Warnln(`no CI detected please use "awesome-ci connect [scm-portal]" to connect!`)
		return
	}

	return
}
