package gitController

var CiEnvironment CIEnvironment

// SetEnvVariables
func SetEnvVariables(envToAppend []string) (err error) {
	switch CiEnvironment.RunnerType {
	case "github_runner":
		err = github_runner_appendToEnv(envToAppend)
	}
	return
}
