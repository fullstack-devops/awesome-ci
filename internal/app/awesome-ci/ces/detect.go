package ces

import (
	"os"

	githubrunner "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces/github_runner"
	gitlabrunner "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces/gitlab_runner"
	jenkinspipeline "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces/jenkins_pipeline"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/models"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"

	log "github.com/sirupsen/logrus"
)

// DetectCes detects the current "code execution side"
// if DetectCes finds an .rc file, this will be used first
func DetectCes() (cesType CES,
	scmPortalType rcpersist.SCMPortalType,
	connCreds *models.ConnectCredentials,
	err error) {

	rcFile := rcpersist.NewRcInstance()
	envFile := "awesomeci.env"

	if creds, errRc := rcFile.Load(); errRc == nil {
		cesType = CES{
			Type:    rcpersist.CESTypeLoMa,
			EnvFile: envFile,
		}
		scmPortalType = rcFile.SCMPortalType
		connCreds = &models.ConnectCredentials{
			ServerURL:  creds.ServerURL,
			Repository: creds.Repository,
			Token:      *creds.TokenPlain,
		}
		return
	}

	// try github runner
	if connCreds, _ = githubrunner.DetectGitHubActionsRunner(); connCreds != nil {
		ghOutput := os.Getenv("GITHUB_OUTPUT")
		cesType = CES{
			Type:    rcpersist.CESTypeGitHubRunner,
			EnvFile: os.Getenv("GITHUB_ENV"),
			OutFile: &ghOutput,
		}
		scmPortalType = rcpersist.SCMPortalTypeGitHub
		return
	}

	// try gitlab runner
	if connCreds, _ = gitlabrunner.DetectGitLabActionsRunner(); connCreds != nil {
		cesType = CES{
			Type:    rcpersist.CESTypeGitLabRunner,
			EnvFile: envFile,
		}
		scmPortalType = rcpersist.SCMPortalTypeGitLab
		return
	}

	// try jenkins pipeline
	if connCreds, _ = jenkinspipeline.DetectJenkinsPipeline(); connCreds != nil {
		cesType = CES{
			Type:    rcpersist.CESTypeJenkinsPipeline,
			EnvFile: envFile,
		}
		scmPortalType = rcpersist.SCMPortalTypeGitHub
		return
	}

	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "" && ciPresent

	if !isCI {
		cesType = CES{
			Type:    rcpersist.CESTypeLoMa,
			EnvFile: envFile,
		}
		log.Warnln(`no CI detected please use "awesome-ci connect [scm-portal]" to connect!`)
		return
	}

	return
}
