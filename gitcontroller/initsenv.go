package gitcontroller

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string, confidential bool) string {
	if value, ok := os.LookupEnv(key); ok {
		var newValue = "**********"
		if len(value) > 10 {
			newValue = value[0:4] + "******"
		}
		return newValue
	}
	return fallback
}

func determanEnvironment(environment string) GitConfiguration {
	var conf GitConfiguration
	switch environment {
	case "Github":
		return GitConfiguration{
			ApiUrl:      os.Getenv("GITHUB_API_URL"),
			Repository:  os.Getenv("GITHUB_REPOSITORY"),
			AccessToken: os.Getenv("GITHUB_TOKEN"),
		}

	case "Custom":
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
	switch environment {
	case "Github":
		conf := determanEnvironment(environment)
		return github_getLatestReleaseVersion(conf)
	}
	return ""
}

// GetLatestReleaseVersion
func CreateNextGitHubRelease(environment string, newReleaseVersion string, uploadArtifacts string) {
	switch environment {
	case "Github":
		conf := determanEnvironment(environment)
		github_createNextGitHubRelease(conf, newReleaseVersion, uploadArtifacts)
	}
}
