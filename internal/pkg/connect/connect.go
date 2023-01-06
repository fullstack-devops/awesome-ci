package connect

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func RemoveRcFile() {
	err := os.Remove(rcFileName) // remove a single file
	if err != nil {
		log.Fatalf("error at removing %s: %v", rcFileName, err)
	}
}

func CheckConnection() {
	rcFile := NewRcInstance()

	if rcFile.Exists() {
		log.Infoln("found existing .rc file")

		creds, err := rcFile.Load()
		if err != nil {
			log.Fatalln(err)
		}

		switch rcFile.Type {
		case serverTypeGitHub:
			_, err = NewGitHubClient(*creds)
			if err != nil {
				log.Fatalf("connection to github could not be established please check. %v", err)
			}
			log.Infof("Successfully connected to github with .rc file to %s", rcFile.ConnectCreds.ServerUrl)

		case serverTypeGitLab:
			log.Warnf("gitlab not yet implemented")

		default:
			log.Fatal("Type in rcFile not known")
		}

	} else {
		log.Fatalf("no %s found, please connect first", rcFileName)
	}
}
