package connect

import (
	"awesome-ci/internal/pkg/githubapi"
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/tools"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ServerType string

const (
	rcFileName                  = ".awesomerc.yaml"
	serverTypeGitHub ServerType = "github"
	serverTypeGitLab ServerType = "gitlab"
)

var (
	rcFile RcFile
	key    = []byte("TZPtSIacEJG18IpqQSkTE6luYmnCNKgR")
)

type RcFile struct {
	Type       ServerType
	Token      string
	Server     string
	Repository string
}

func RemoveRcFile() {
	err := os.Remove(rcFileName) // remove a single file
	if err != nil {
		log.Fatalf("error at removing %s: %v", rcFileName, err)
	}
}

func CheckConnection() {
	if tools.CheckFileExists(rcFileName) {
		log.Infof("found existing %s file", rcFileName)
		yamlFile, err := os.ReadFile(rcFileName)
		if err != nil {
			log.Fatalf("Error while opening existing rc file. %v", err)
		}
		log.Tracef("Successfully read file %s", rcFileName)

		err = yaml.Unmarshal(yamlFile, &rcFile)
		if err != nil {
			log.Fatalf("error at reading file %s: %v", rcFileName, err)
		}
		log.Tracef("Successfully unmarshal file %s", rcFileName)

		decryptedToken, err := tools.DecryptString(rcFile.Token, key)
		if err != nil {
			log.Fatal(err)
		}
		log.Tracef("Successfully decrypt string (github token %s***)", decryptedToken[0:8])

		// check if connection can be established
		switch rcFile.Type {
		case serverTypeGitHub:
			_, err = githubapi.NewGitHubClient(models.StandardConnectCredentials{
				ServerUrl:  rcFile.Server,
				Repository: rcFile.Repository,
				Token:      string(decryptedToken),
			})
			if err != nil {
				log.Fatalf("connection to github could not be established please check. %v", err)
			}
			log.Infof("Successfully connected to github %s", rcFile.Server)

		case serverTypeGitLab:
			log.Warnf("gitlab not yet implemented")

		default:
			log.Fatal("Type in rcFile not known")
		}

	} else {
		log.Fatalf("no %s found, please connect first", rcFileName)
	}
}
