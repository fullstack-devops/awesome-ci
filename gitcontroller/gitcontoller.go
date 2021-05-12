package gitcontroller

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getEnv(key, fallback string, confidential bool) string {
	if value, ok := os.LookupEnv(key); ok {
		var newValue = "**********"
		if confidential {
			if len(value) > 10 {
				newValue = value[0:4] + "******"
			}
		} else {
			newValue = value
		}
		return newValue
	}
	return fallback
}

func evaluateEnvironment(overirdecienv string) (cienv string) {
	if overirdecienv != "" {
		if strings.Contains("github_runner|jenkins|gitlab_ci", overirdecienv) {
			cienv = overirdecienv
		} else {
			log.Fatalf("%s is not a valide cienv\nValid are: (github_runner|jenkins|gitlab_ci)", overirdecienv)
		}
	} else {
		// check if github-runner
		_, githubRunnerApi := os.LookupEnv("GITHUB_API_URL")
		_, githubRunnerRep := os.LookupEnv("GITHUB_REPOSITORY")
		if githubRunnerApi && githubRunnerRep {
			_, githubRunnerToken := os.LookupEnv("GITHUB_TOKEN")
			if !githubRunnerToken {
				log.Fatalln("Apparently you are using a GitHub-Runner.\nPlease provide the GITHUB_TOKEN!\nSee ")
			}
			return "github_runner"
		}

	}

	switch cienv {
	case "github_runner":
		_, githubRunnerToken := os.LookupEnv("GITHUB_TOKEN")
		if !githubRunnerToken {
			log.Fatalln("Apparently you are using a GitHub-Runner.\nPlease provide the GITHUB_TOKEN!\nSee https://eksrvb.github.io/awesome-ci/examples/github_actions.html")
		}
		return "github_runner"
	default:
		return "github_runner"
	}
}

func determanEnvironment(environment string) GitConfiguration {
	var conf GitConfiguration
	switch environment {
	case "github_runner":
		return GitConfiguration{
			ApiUrl:            os.Getenv("GITHUB_API_URL"),
			Repository:        os.Getenv("GITHUB_REPOSITORY"),
			AccessToken:       os.Getenv("GITHUB_TOKEN"),
			DefaultBranchName: getEnv("GIT_DEFAULT_BRANCH_NAME", "main", false),
		}

	case "custom":
		fmt.Println("In order to establish a connection to the github please provide the following environment variables:")
		fmt.Printf("   GIT_HOSTNAME      %s\n", getEnv("GIT_HOSTNAME", "git.daimler.com", false))
		fmt.Printf("   GIT_API_VERSION   %s\n", getEnv("GIT_API_VERSION", "v3", false))
		fmt.Printf("   GIT_ORGA          %s\n", getEnv("GIT_ORGA", "not-set", false))
		fmt.Printf("   GIT_REPO          %s\n", getEnv("GIT_REPO", "not-set", false))
		fmt.Printf("   GIT_ACCESS_TOKEN  %s\n", getEnv("GIT_ACCESS_TOKEN", "not-set", true))
	}

	return conf
}

// GetLatestReleaseVersion
func GetLatestReleaseVersion(environment string) string {
	switch evaluateEnvironment(environment) {
	case "github_runner":
		conf := determanEnvironment(environment)
		return github_getLatestReleaseVersion(conf)
	}
	return ""
}

// GetLatestReleaseVersion
func CreateNextGitHubRelease(environment string, newReleaseVersion string, uploadArtifacts string) {
	switch evaluateEnvironment(environment) {
	case "github_runner":
		conf := determanEnvironment(environment)
		github_createNextGitHubRelease(conf, newReleaseVersion, uploadArtifacts)
	}
}
