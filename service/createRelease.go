package service

import (
	"awesome-ci/gitcontroller"
	"awesome-ci/semver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CreateRelease(cienv string, overrideVersion *string, getVersionIncrease *string, isDryRun *bool, publishNpm *string, uploadArtifacts *string) {
	var gitVersion string
	if *overrideVersion != "" {
		gitVersion = *overrideVersion
	} else {
		gitVersion = gitcontroller.GetLatestReleaseVersion(cienv)
	}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		if cienv == "Github" {
			buildInfos, err := getMergeMessage()
			if err != nil {
				log.Fatal(err)
			} else {
				patchLevel = buildInfos.PatchLevel
			}
		}
	}

	newVersion := semver.IncreaseSemVer(patchLevel, gitVersion)
	if *isDryRun {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Would writing new release: %s\n", newVersion)
	} else {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Writing new release: %s\n", newVersion)
		gitcontroller.CreateNextGitHubRelease(cienv, newVersion, *uploadArtifacts)
	}

	if *publishNpm != "" {
		// check if subfolder has slash
		pathToSource := *publishNpm
		if !strings.HasSuffix(*publishNpm, "/") {
			pathToSource = *publishNpm + "/"
		}
		fmt.Printf("Puplishing npm packages under path: %s\n", pathToSource)
		npmPublish(pathToSource, newVersion)
	}
}

func npmPublish(pathToSource string, newVersion string) {

	// opening package.json
	jsonFile, err := os.Open(pathToSource + "package.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer jsonFile.Close()

	var result map[string]interface{}
	json.NewDecoder(jsonFile).Decode(&result)

	result["version"] = newVersion

	b, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// writing result to package.json
	err = ioutil.WriteFile(pathToSource+"package.json", b, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	command := exec.Command("npm", "publish", pathToSource, "--tag", fmt.Sprintf("%s@%s", result["name"], newVersion))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// Run the command
	command.Run()
}
