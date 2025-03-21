package connect

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"

	log "github.com/sirupsen/logrus"
)

func CheckConnection() {
	rcFile := rcpersist.NewRcInstance()

	if rcFile.Exists() {
		log.Infoln("found existing .rc file")

		creds, err := rcFile.Load()
		if err != nil {
			log.Fatalln(err)
		}

		switch rcFile.SCMPortalType {
		case rcpersist.SCMPortalTypeGitHub:
			_, err = github.NewGitHubClient(&creds.ServerURL, &creds.Repository, creds.TokenPlain)
			if err != nil {
				log.Fatalf("connection to github could not be established please check. %v", err)
			}
			log.Infof("Successfully connected to github with .rc file to %s", rcFile.ConnectCreds.ServerURL)

		case rcpersist.SCMPortalTypeGitLab:
			log.Warnf("gitlab not yet implemented")

		default:
			log.Fatal("Type in rcFile not known")
		}

	} else {
		log.Fatalf("no .rc found, please connect first")
	}
}
