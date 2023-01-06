package detect

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func LoadEnvVars() (envs EnvVariables, err error) {
	// in github actions the en file should be used. We use this to detect a github actions environment
	_, isgithubEnv := os.LookupEnv("GITHUB_ENV")
	if isgithubEnv {
		log.Traceln("loading envs as github actions runner from file")
		envs.CiType = GitHubActionsRunner
		envs.Envs, err = openGitHubActionsEnvFile()
		return
	} else {
		log.Traceln("loading envs directly from os")
		allEnvs := os.Environ()

		envs.CiType = NormalEnv
		for _, v := range allEnvs {
			envs.Envs = append(envs.Envs, devideEnvStringToKeyAndValue(v))
		}
	}
	return
}

func (envs *EnvVariables) SetEnvVars() (err error) {
	if envs.CiType == GitHubActionsRunner {
		return saveToGitHubActionsEnvFile(envs.Envs)
	} else {
		/* for _, env := range envs.Envs {
			os.Unsetenv()
			os.Setenv(*env.Name, *env.Value)
		} */
		return nil
	}
}
