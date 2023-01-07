package detect

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func openGitHubActionsEnvFile() (variables []EnvVariable, err error) {
	envFile, err := os.Open(os.Getenv("GITHUB_ENV"))
	if errors.Is(err, os.ErrNotExist) {
		envFile, err = os.Create(os.Getenv("GITHUB_ENV"))
	}

	defer envFile.Close()

	if err != nil {
		return nil, fmt.Errorf("error at opening env file: %v", err)
	}

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		s := scanner.Text()
		variables = append(variables, devideEnvStringToKeyAndValue(s))
	}
	return
}

func saveToGitHubActionsEnvFile(variables []EnvVariable) (err error) {
	envFile, err := os.OpenFile(os.Getenv("GITHUB_ENV"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = fmt.Errorf("error at opening env file: %v", err)
		return
	}

	defer envFile.Close()

	envFile.Truncate(0)
	wirteEnvs := ""
	for _, env := range variables {
		wirteEnvs = wirteEnvs + fmt.Sprintf("%s\n", env.ToString())
	}
	_, err = envFile.Write([]byte(wirteEnvs))
	if err != nil {
		return
	}
	err = envFile.Close()
	return
}
