package service

import (
	"awesome-ci/gitOnlineController"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CreateRelease(cienv string, overrideVersion *string, getVersionIncrease *string, isDryRun *bool, preRelease *bool, publishNpm *string, uploadArtifacts *string) {
	var gitVersion string
	if *overrideVersion != "" {
		gitVersion = *overrideVersion
	} else {
		gitVersion = gitOnlineController.GetLatestReleaseVersion()
	}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		buildInfos, err := getLatestCommitMessage()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			patchLevel = buildInfos.PatchLevel
		}
	}

	newVersion := increaseSemVer(patchLevel, gitVersion)
	if *isDryRun {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Would writing new release: %s\n", newVersion)
	} else {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Writing new release: %s\n", newVersion)

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
	gitOnlineController.CreateNextGitHubRelease(CiEnvironment.GitInfos.DefaultBranchName, newVersion, preRelease, isDryRun, uploadArtifacts)
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
