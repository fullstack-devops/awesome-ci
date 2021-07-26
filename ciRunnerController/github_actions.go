package ciRunnerController

import (
	"awesome-ci/models"
	"fmt"
	"log"
	"os"
)

func github_runner_appendToEnv(envToAppend []models.BuildEnvironmentVariable) (err error) {
	f, err := os.OpenFile(CiEnvironment.RunnerInfo.EnvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error at opening GITHUB_ENV file:", err)
	}

	var dstring string = ""
	for _, env := range envToAppend {
		dstring = fmt.Sprintf("%s\n%s=%s", dstring, env.Key, env.Value)
	}

	if _, err := f.Write([]byte(dstring)); err != nil {
		log.Fatal("Error at writing to GITHUB_ENV file:", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal("Error at closing GITHUB_ENV file:", err)
	}
	return
}
