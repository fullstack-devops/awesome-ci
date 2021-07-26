package ciRunnerController

import "awesome-ci/models"

var CiEnvironment models.CIEnvironment

// SetEnvVariables
func SetEnvVariables(envToAppend []models.BuildEnvironmentVariable) (err error) {
	switch CiEnvironment.RunnerType {
	case "github_runner":
		err = github_runner_appendToEnv(envToAppend)
	}
	return
}
