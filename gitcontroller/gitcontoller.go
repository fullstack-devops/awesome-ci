package gitcontroller

import (
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

	// goes only at cienv override
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

// GetLatestReleaseVersion
func GetLatestReleaseVersion(environment string) string {
	switch evaluateEnvironment(environment) {
	case "github_runner":
		return github_getLatestReleaseVersion(
			os.Getenv("GITHUB_API_URL"),
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("GITHUB_TOKEN"))
	}
	return ""
}

// GetLatestReleaseVersion
func CreateNextGitHubRelease(environment string, releaseBranch string, newReleaseVersion string, preRelease *bool, uploadArtifacts string) {
	switch evaluateEnvironment(environment) {
	case "github_runner":
		github_createNextGitHubRelease(
			os.Getenv("GITHUB_API_URL"),
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("GITHUB_TOKEN"),
			releaseBranch,
			newReleaseVersion,
			*preRelease,
			uploadArtifacts)
	}
}
