package connect

import (
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/pkg/rcpersist"

	log "github.com/sirupsen/logrus"
)

func UpdateRcFileForGitHub(server string, repo string, token string) {
	rcFile := rcpersist.NewRcInstance()

	_, err := rcFile.Load()
	if err != nil && err != rcpersist.ErrRcFileNotExists {
		log.Fatalln(err)
	}

	rcFile.UpdateSCMPortalType(rcpersist.SCMPortalTypeGitHub)
	rcFile.UpdateCreds(server, repo, token)

	if err := rcFile.Save(); err != nil {
		log.Fatalln(err)
	}

	// testing connection
	if _, err := github.NewGitHubClient(&server, &repo, &token); err != nil {
		log.Fatalln(err)
	}
}
