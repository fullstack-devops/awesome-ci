package acigithub

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	envVars                EnvVariables
	githubEnv, isgithubEnv = os.LookupEnv("GITHUB_ENV")
)

type EnvVariables struct {
	envs     []EnvVariable
	fileName *string
}

type EnvVariable struct {
	Name  *string
	Value *string
}

func OpenEnvFile() (envVariables *EnvVariables, err error) {
	if isgithubEnv {
		envFile, err := os.Open(githubEnv)

		if errors.Is(err, os.ErrNotExist) {
			envFile, err = os.Create(githubEnv)
		}

		if err != nil {
			return nil, errors.New(fmt.Sprintln("Error at opening ENV file:", err))
		}
		envVars.fileName = &githubEnv

		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			s := scanner.Text()
			envVars.Set(s[:strings.Index(s, "=")], s[strings.Index(s, "=")+1:])
		}
	}
	return &envVars, nil
}

func (envs *EnvVariables) SaveEnvFile() (err error) {
	envFile, err := os.OpenFile(*envs.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = errors.New(fmt.Sprintln("Error at opening ENV file:", err))
		return
	}
	envFile.Truncate(0)
	wirteEnvs := ""
	for _, env := range envs.envs {
		wirteEnvs = wirteEnvs + fmt.Sprintf("%s=%s\n", *env.Name, *env.Value)
	}
	_, err = envFile.Write([]byte(wirteEnvs))
	if err != nil {
		return
	}
	err = envFile.Close()
	return
}

func (envs *EnvVariables) Set(name string, value string) {
	var envFound bool = false
	for _, env := range envs.envs {
		if *env.Name == name {
			*env.Value = value
			envFound = true
			return
		}
	}
	if !envFound {
		envs.envs = append(envs.envs, EnvVariable{Name: &name, Value: &value})
	}
}

func (envs *EnvVariables) Get(name string) (envVariable *EnvVariable) {
	for _, env := range envs.envs {
		if *env.Name == name {
			return &env
		}
	}
	return nil
}
