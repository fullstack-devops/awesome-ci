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

func CreateRelease(cienv string, versionOverr *string, patchLevelOverr *string, isDryRun *bool, preRelease *bool, publishNpm *string, uploadArtifacts *string) {

	_, _, buildInfos := setBuildInfos(versionOverr, patchLevelOverr)

	if *isDryRun {
		fmt.Printf("Old version: %s\n", buildInfos.Version)
		fmt.Printf("Would writing new release: %s\n", buildInfos.NextVersion)
	} else {
		fmt.Printf("Old version: %s\n", buildInfos.Version)
		fmt.Printf("Writing new release: %s\n", buildInfos.NextVersion)

		if *publishNpm != "" {
			// check if subfolder has slash
			pathToSource := *publishNpm
			if !strings.HasSuffix(*publishNpm, "/") {
				pathToSource = *publishNpm + "/"
			}
			fmt.Printf("Puplishing npm packages under path: %s\n", pathToSource)
			npmPublish(pathToSource, buildInfos.NextVersion)
		}
	}
	gitOnlineController.CreateNextGitHubRelease(CiEnvironment.GitInfos.DefaultBranchName, buildInfos.NextVersion, preRelease, isDryRun, uploadArtifacts)
}

func npmPublish(pathToSource string, nextVersion string) {

	// opening package.json
	jsonFile, err := os.Open(pathToSource + "package.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer jsonFile.Close()

	var result map[string]interface{}
	json.NewDecoder(jsonFile).Decode(&result)

	result["version"] = nextVersion

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

	command := exec.Command("npm", "publish", pathToSource, "--tag", fmt.Sprintf("%s@%s", result["name"], nextVersion))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// Run the command
	command.Run()
}
