package service

import (
	"awesome-ci/src/semver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v39/github"
)

func ReleaseCreate(versionOverr *string, patchLevelOverr *string, isDryRun *bool) {

	prInfos, prNumber, err := getPRInfos(nil)
	if err != nil {
		panic(err)
	}

	branchName := *prInfos.Head.Ref
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
		repositoryRelease, err := CiEnvironment.GetLatestReleaseVersion()
		if err != nil {
			log.Println(err)
		}
		gitVersion = *repositoryRelease.TagName
	}
	nextVersion, err := semver.IncreaseVersion(patchLevel, gitVersion)

	if *isDryRun {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Would writing new release: %s\n", nextVersion)
	} else {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Writing new release: %s\n", nextVersion)
		relName := "Release " + nextVersion
		draft := true

		releaseObject := github.RepositoryRelease{
			TargetCommitish: &branchName,
			TagName:         &nextVersion,
			Name:            &relName,
			Draft:           &draft,
		}

		CiEnvironment.ManageGitRelease(releaseObject, nil)
	}
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
