package connect

import (
	"awesome-ci/internal/pkg/githubapi"
	"awesome-ci/internal/pkg/tools"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ConnectToGitHub() (githubRichClient *githubapi.GitHubRichClient, err error) {
	if tools.CheckFileExists(rcFileName) {
		log.Info("Connecting to Github via local %s file", rcFileName)
		return ConnectWithRcFileToGithub()
	} else {
		isGithubActionsRunner, _, _ := CheckGitHubRunner()
		if isGithubActionsRunner {
			log.Info("Connecting to Github via GitHub actions runner")
			return ConnectWithActionsToGitHub()
		}
		isJenkinsPipeline, _, _ := CheckJenkinkPipeline()
		if isJenkinsPipeline {
			log.Info("Connecting to Github via Jenkins pipeline")
			return ConnectWithJenkinsToGitHub()
		}
	}
	return nil, fmt.Errorf("could not connect with any method with GitHub. Please read the Docs")
}

func UpdateRcFileForGithub(server string, repo string, token string) {
	if tools.CheckFileExists(rcFileName) {
		log.Tracef("found existing %s file, try to update...", rcFileName)
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

		// check if connection can be established
		decryptedToken, err := tools.DecryptString(rcFile.Token, key)
		if err != nil {
			log.Fatal(err)
		}
		log.Tracef("Successfully decrypt string (github token %s)", decryptedToken[0:8])

		_, err = githubapi.NewGitHubClient(githubapi.ConnectCredentials{
			ServerUrl:  rcFile.Server,
			Repository: rcFile.Repository,
			Token:      string(decryptedToken),
		})
		if err != nil {
			log.Fatalf("connection to github could not be established please check. %v", err)
		}
		log.Tracef("Successfully connected to github %s", rcFile.Server)

	} else {

		encryptedToken, err := tools.EncryptString(token, key)
		if err != nil {
			log.Error(err)
		}
		rcFile = RcFile{
			Server:     server,
			Repository: repo,
			Token:      string(encryptedToken),
			Type:       serverTypeGitHub,
		}

		// check connection
		_, err = githubapi.NewGitHubClient(githubapi.ConnectCredentials{
			ServerUrl:  rcFile.Server,
			Repository: rcFile.Repository,
			Token:      string(encryptedToken),
		})
		if err != nil {
			log.Fatalf("connection to github could not be established please check. %v", err)
		}
	}

	yamlData, err := yaml.Marshal(&rcFile)
	if err != nil {
		log.Fatalf("Error while Marshaling. %v", err)
	}

	err = os.WriteFile(rcFileName, yamlData, 0600)
	if err != nil {
		log.Fatalf("Unable to write data into rc file")
	}
}

func ConnectWithRcFileToGithub() (githubRichClient *githubapi.GitHubRichClient, err error) {
	if tools.CheckFileExists(rcFileName) {
		yamlFile, err := os.ReadFile(rcFileName)
		if err != nil {
			log.Fatalf("Error while opening existing rc file. %v", err)
		}
		log.Tracef("Successfully read file %s", rcFileName)

		err = yaml.Unmarshal(yamlFile, &rcFile)
		if err != nil {
			log.Fatalf("error at reading file %s: %v", rcFileName, err)
		}

		// check if connection can be established
		decryptedToken, err := tools.DecryptString(rcFile.Token, key)
		if err != nil {
			log.Fatal(err)
		}

		return githubapi.NewGitHubClient(githubapi.ConnectCredentials{
			ServerUrl:  rcFile.Server,
			Repository: rcFile.Repository,
			Token:      string(decryptedToken),
		})

	} else {

		return nil, fmt.Errorf("no %s file found", rcFileName)
	}
}

// GitHub Actions
func ConnectWithActionsToGitHub() (githubRichClient *githubapi.GitHubRichClient, err error) {
	isRunner, creds, _ := CheckGitHubRunner()
	if isRunner {
		return githubapi.NewGitHubClient(creds)
	} else {
		err = fmt.Errorf("")
		return
	}
}
func CheckGitHubRunner() (isRunner bool, connects githubapi.ConnectCredentials, err error) {
	githubServerUrl, isgithubServerUrl := os.LookupEnv("GITHUB_SERVER_URL")
	githubRepository, isgithubRepository := os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken := os.LookupEnv("GITHUB_TOKEN")

	if isgithubServerUrl && isgithubRepository && isgithubToken {
		return true, githubapi.ConnectCredentials{
			ServerUrl:  githubServerUrl,
			Repository: githubRepository,
			Token:      githubToken,
		}, nil

	} else {
		return false, githubapi.ConnectCredentials{}, fmt.Errorf("no github actions runner detected")
	}
}

// GitHub Actions
func ConnectWithJenkinsToGitHub() (githubRichClient *githubapi.GitHubRichClient, err error) {
	isRunner, creds, _ := CheckJenkinkPipeline()
	if isRunner {
		return githubapi.NewGitHubClient(creds)
	} else {
		err = fmt.Errorf("")
		return
	}
}
func CheckJenkinkPipeline() (isJenkins bool, connects githubapi.ConnectCredentials, err error) {
	_, isjenkinsUrl := os.LookupEnv("JENKINS_URL")
	githubServerUrl, isgithubServerUrl := os.LookupEnv("GIT_URL")
	githubRepository, isgithubRepository := os.LookupEnv("GITHUB_REPOSITORY")
	githubToken, isgithubToken := os.LookupEnv("GITHUB_TOKEN")

	if isjenkinsUrl && isgithubServerUrl && isgithubRepository && isgithubToken {
		return true, githubapi.ConnectCredentials{
			ServerUrl:  githubServerUrl,
			Repository: githubRepository,
			Token:      githubToken,
		}, nil

	} else {
		return false, githubapi.ConnectCredentials{}, fmt.Errorf("no jenkins pipeline detected")
	}
}
