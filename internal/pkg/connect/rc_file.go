package connect

import (
	"awesome-ci/internal/pkg/tools"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	rcTokenKey = []byte("TZPtSIacEJG18IpqQSkTE6luYmnCNKgR")
)

const (
	rcFileName                  = ".awesomerc.yaml"
	serverTypeGitHub ServerType = "github"
	serverTypeGitLab ServerType = "gitlab"
)

func (t Token) ToString() string {
	return string(t)
}

func (t Token) Decrypt() (string, error) {
	if decryptedToken, err := tools.DecryptString(t.ToString(), rcTokenKey); err != nil {
		return "", fmt.Errorf("could not decrypt token: %v", err)
	} else {
		return decryptedToken, nil
	}
}

func (t Token) Encrypt() (string, error) {
	if encryptedToken, err := tools.EncryptString(t.ToString(), rcTokenKey); err != nil {
		return "", fmt.Errorf("could not encrypt token: %v", err)
	} else {
		return encryptedToken, nil
	}
}

func NewRcInstance() RcFile {
	return RcFile{}
}

func (rc RcFile) Exists() bool {
	return tools.CheckFileExists(rcFileName)
}

func (rc RcFile) Load() (creds *ConnectCredentials, err error) {
	if !rc.Exists() {
		return nil, fmt.Errorf("rc file does not exist. please connect first")
	}

	// load rc file
	yamlFile, err := os.ReadFile(rcFileName)
	if err != nil {
		return nil, fmt.Errorf("Error while opening existing rc file. %v", err)
	}
	log.Tracef("Successfully read file %s", rcFileName)

	if err = yaml.Unmarshal(yamlFile, &rc); err != nil {
		return nil, fmt.Errorf("error at reading file %s: %v", rcFileName, err)
	}
	log.Tracef("Successfully unmarshal file %s", rcFileName)

	if rc.ConnectCreds.TokenPlain, err = rc.ConnectCreds.Token.Encrypt(); err != nil {
		return nil, err
	}
	return
}

func (rc RcFile) UpdateServerType(serverType ServerType) {
	rc.Type = serverType
}

func (rc RcFile) UpdateCreds(server string, repo string, token string) error {
	if rc.ConnectCreds.Token.ToString() != token {
		rc.ConnectCreds.Token = Token(token)
		if decryptedToken, err := Token(token).Decrypt(); err != nil {
			return err
		} else {
			rc.ConnectCreds.TokenPlain = decryptedToken
		}
	}
	if rc.ConnectCreds.ServerUrl != server {
		rc.ConnectCreds.ServerUrl = server
	}
	if rc.ConnectCreds.Repository != repo {
		rc.ConnectCreds.Repository = repo
	}
	return nil
}

func (rc RcFile) Save() (err error) {
	rc.ConnectCreds.TokenPlain = ""

	os.Truncate(rcFileName, 0)
	yamlData, err := yaml.Marshal(rc)
	if err != nil {
		return fmt.Errorf("Error while Marshaling. %v", err)
	}
	if err := os.WriteFile(rcFileName, yamlData, 0600); err != nil {
		return fmt.Errorf("Unable to write data into rc file")
	}
	log.Tracef("Successfully written config to %s", rcFileName)

	return nil
}
