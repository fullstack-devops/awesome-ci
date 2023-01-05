package detect

import (
	"fmt"
	"os"
	"strings"

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

func (envs *EnvVariables) Set(name string, value string) {
	var envFound bool = false
	for _, env := range envs.Envs {
		if *env.Name == name {
			*env.Value = value
			envFound = true
			return
		}
	}
	if !envFound {
		envs.Envs = append(envs.Envs, EnvVariable{Name: &name, Value: &value})
	}
	os.Setenv(name, value)
}

func (envs *EnvVariables) Get(name string) (envVariable *EnvVariable) {
	for _, env := range envs.Envs {
		if *env.Name == name {
			return &env
		}
	}
	return nil
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

func devideEnvStringToKeyAndValue(envString string) (env EnvVariable) {
	name := envString[:strings.Index(envString, "=")]
	value := envString[strings.Index(envString, "=")+1:]
	return EnvVariable{Name: &name, Value: &value}
}

func (envVar EnvVariable) ToString() (envString string) {
	return fmt.Sprintf("%s=%s", *envVar.Name, *envVar.Value)
}
