package rcpersist

import (
	"fmt"
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

func NewRcInstance() *RcFile {
	return &RcFile{}
}

func (rc *RcFile) Exists() bool {
	return tools.CheckFileExists(rcFileName)
}

func (rc *RcFile) Load() (creds *ConnectCredentials, err error) {
	if !rc.Exists() {
		return nil, ErrRcFileNotExists
	}

	// load rc file
	yamlFile, err := os.ReadFile(rcFileName)
	if err != nil {
		return nil, fmt.Errorf("error while opening existing rc file. %v", err)
	}
	log.Tracef("Successfully read file %s", rcFileName)

	if err = yaml.Unmarshal(yamlFile, &rc); err != nil {
		return nil, fmt.Errorf("error at reading file %s: %v", rcFileName, err)
	}
	log.Tracef("Successfully unmarshal file %s", rcFileName)

	plainToken, err := rc.ConnectCreds.Token.Decrypt()
	if err != nil {
		return nil, err
	}
	rc.ConnectCreds.TokenPlain = &plainToken
	return &rc.ConnectCreds, nil
}

func (rc *RcFile) UpdateSCMPortalType(scmPortalType SCMPortalType) {
	rc.SCMPortalType = scmPortalType
}

func (rc *RcFile) UpdateCESType(cesType CESType) {
	rc.CESType = cesType
}

func (rc *RcFile) UpdateCreds(server string, repo string, token string) error {
	if rc.ConnectCreds.Token.ToString() != token {
		rc.ConnectCreds.TokenPlain = &token
		if encryptedToken, err := Token(token).Encrypt(); err != nil {
			return err
		} else {
			rc.ConnectCreds.Token = Token(encryptedToken)
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

func (rc *RcFile) Save() (err error) {
	rc.ConnectCreds.TokenPlain = nil
	os.Truncate(rcFileName, 0)

	yamlData, err := yaml.Marshal(rc)
	if err != nil {
		return fmt.Errorf("error while Marshaling. %v", err)
	}
	if err := os.WriteFile(rcFileName, yamlData, 0600); err != nil {
		return fmt.Errorf("unable to write data into rc file")
	}
	log.Tracef("Successfully written config to %s", rcFileName)

	return checkGitIgnore()
}

func RemoveRcFile() {
	err := os.Remove(rcFileName) // remove a single file
	if err != nil {
		log.Fatalf("error at removing %s: %v", rcFileName, err)
	}
}
