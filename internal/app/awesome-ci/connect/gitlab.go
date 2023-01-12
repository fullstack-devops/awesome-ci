package connect

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"

	log "github.com/sirupsen/logrus"
)

func UpdateRcFileForGitLab(server string, repo string, token string) {
	rcFile := rcpersist.NewRcInstance()

	_, err := rcFile.Load()
	if err != nil {
		log.Fatalln(err)
	}

	rcFile.UpdateSCMPortalType(rcpersist.SCMPortalTypeGitLab)
	rcFile.UpdateCreds(server, repo, token)

	if err := rcFile.Save(); err != nil {
		log.Fatalln(err)
	}

	// testing connection
	if _, err := gitlab.NewGitLabClient(&server, &repo, &token); err != nil {
		log.Fatalln(err)
	}
}
