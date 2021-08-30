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

	prInfos, prNumber, err := getPRInfos()
	if err != nil {
		panic(err)
	}

	branchName := prInfos.Head.Ref
	patchLevel := branchName[:strings.Index(branchName, "/")]

	if *patchLevelOverr != "" {
		patchLevel = *patchLevelOverr
	}

	// if an comment exists with aci=major, make a major version!
	if detectIfMajor(prNumber) {
		patchLevel = "major"
	}

	var gitVersion string
	if *versionOverr != "" {
		gitVersion = *versionOverr
	} else {
		gitVersion = gitOnlineController.GetLatestReleaseVersion()
	}
	nextVersion := increaseSemVer(patchLevel, gitVersion)

	if *isDryRun {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Would writing new release: %s\n", nextVersion)
	} else {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Writing new release: %s\n", nextVersion)

		if *publishNpm != "" {
			// check if subfolder has slash
			pathToSource := *publishNpm
			if !strings.HasSuffix(*publishNpm, "/") {
				pathToSource = *publishNpm + "/"
			}
			fmt.Printf("Puplishing npm packages under path: %s\n", pathToSource)
			npmPublish(pathToSource, nextVersion)
		}
	}
	gitOnlineController.CreateNextGitHubRelease(CiEnvironment.GitInfos.DefaultBranchName, nextVersion, preRelease, isDryRun, uploadArtifacts)
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

	command := exec.Command("npm", "publish", pathToSource, "--tag latest")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// Run the command
	command.Run()
}
