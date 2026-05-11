package connect

import (
	"github.com/layer87-labs/relctl/internal/app/relctl/scm-portal/github"
	"github.com/layer87-labs/relctl/internal/pkg/rcpersist"

	log "github.com/sirupsen/logrus"
)

func UpdateRcFileForGitHub(server string, repo string, token string) {
	rcFile := rcpersist.NewRcInstance()

	_, err := rcFile.Load()
	if err != nil && err != rcpersist.ErrRcFileNotExists {
		log.Fatalln(err)
	}

	rcFile.UpdateSCMPortalType(rcpersist.SCMPortalTypeGitHub)
	if err := rcFile.UpdateCreds(server, repo, token); err != nil {
		log.Fatalln(err)
	}

	if err := rcFile.Save(); err != nil {
		log.Fatalln(err)
	}

	// testing connection
	if _, err := github.NewGitHubClient(&server, &repo, &token); err != nil {
		log.Fatalln(err)
	}
	log.Infof("successfully connected to github at %s", server)
}
