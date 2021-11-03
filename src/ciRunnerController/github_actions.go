package ciRunnerController

import (
	"log"
	"os"
)

func github_runner_appendToEnv(envToAppend []string) (err error) {
	f, err := os.OpenFile(CiEnvironment.RunnerInfo.EnvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error at opening GITHUB_ENV file:", err)
	}

	var dstring string = ""
	for _, singleString := range envToAppend {
		dstring = dstring + "\n" + singleString
	}

	if _, err := f.Write([]byte(dstring)); err != nil {
		log.Fatal("Error at writing to GITHUB_ENV file:", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal("Error at closing GITHUB_ENV file:", err)
	}
	return
}
