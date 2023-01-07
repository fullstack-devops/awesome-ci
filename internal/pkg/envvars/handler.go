package envvars

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func OpenEnvFile(otherEnvFile string) (envVars *EnvVariables, err error) {
	envFile, err := os.Open(checkForEnvFileOverride(otherEnvFile))
	// ignore and create at write
	if !errors.Is(err, os.ErrNotExist) {
		return
	}

	defer envFile.Close()

	if err != nil {
		return nil, fmt.Errorf("error at opening env file %v", err)
	}

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		s := scanner.Text()
		envVars.Envs = append(envVars.Envs, devideEnvStringToKeyAndValue(s))
	}
	return
}

func (envVars *EnvVariables) CloseEnvFile(otherEnvFile string) (err error) {
	envFile, err := os.OpenFile(checkForEnvFileOverride(otherEnvFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error at opening env file: %v", err)
	}

	defer envFile.Close()

	envFile.Truncate(0)
	wirteEnvs := ""
	for _, env := range envVars.Envs {
		wirteEnvs = wirteEnvs + fmt.Sprintf("%s\n", env.ToString())
	}

	if _, err = envFile.Write([]byte(wirteEnvs)); err != nil {
		return err
	}

	return envFile.Close()
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
}

func (envs *EnvVariables) Get(name string) (envVariable *EnvVariable) {
	for _, env := range envs.Envs {
		if *env.Name == name {
			return &env
		}
	}
	return nil
}

func devideEnvStringToKeyAndValue(envString string) (env EnvVariable) {
	name := envString[:strings.Index(envString, "=")]
	value := envString[strings.Index(envString, "=")+1:]
	return EnvVariable{Name: &name, Value: &value}
}

func (envVar EnvVariable) ToString() (envString string) {
	return fmt.Sprintf("%s=%s", *envVar.Name, *envVar.Value)
}

func checkForEnvFileOverride(override string) string {
	if override != "" {
		return override
	} else {
		return EnvFile
	}
}
