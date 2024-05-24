package envvars

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func OpenEnvFile(file string) (envVars *EnvVariables, err error) {
	envFile, err := os.Open(file)
	// ignore and create at write
	if errors.Is(err, os.ErrNotExist) {
		return &EnvVariables{}, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error at opening env file %v", err)
	}
	defer envFile.Close()

	envVars = &EnvVariables{}

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		s := scanner.Text()
		envVars.Envs = append(envVars.Envs, devideEnvStringToKeyAndValue(s))
	}
	return
}

func (envVars *EnvVariables) CloseEnvFile(file string) (err error) {
	envFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (envVars *EnvVariables) Set(name string, value string) {
	var envFound = false
	for _, env := range envVars.Envs {
		if *env.Name == name {
			*env.Value = value
			envFound = true
			return
		}
	}
	if !envFound {
		envVars.Envs = append(envVars.Envs, EnvVariable{Name: &name, Value: &value})
	}
}

func (envVars *EnvVariables) Get(name string) (envVariable *EnvVariable) {
	for _, env := range envVars.Envs {
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
