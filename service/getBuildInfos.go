package service

import (
	"awesome-ci/gitcontroller"
	"fmt"
	"strings"
)

func GetBuildInfos(cienv string, overrideVersion *string, getVersionIncrease *string, format *string) {

	var infosMergeMessage infosMergeMessage
	//if cienv == "Github" {
	//	var err error
	infosMergeMessage, err := getLatestCommitMessage()
	if err != nil {
		infosMergeMessage.PRNumber = fmt.Sprint(gitcontroller.GetPrNumberForBranch(getCurrentBranchName()))
	}
	//}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		patchLevel = infosMergeMessage.PatchLevel
	}

	var gitVersion string
	var nextVersion string
	if strings.Contains(*format, "version") || *format == "" {
		if *overrideVersion != "" {
			gitVersion = *overrideVersion
		} else {
			gitVersion = gitcontroller.GetLatestReleaseVersion(cienv)
		}
		nextVersion = increaseSemVer(patchLevel, gitVersion)
	}

	if *format != "" {
		replacer := strings.NewReplacer(
			"pr", infosMergeMessage.PRNumber,
			"version", gitVersion,
			"next_version", nextVersion,
			"patchLevel", patchLevel)
		output := replacer.Replace(*format)
		fmt.Print(output)
	} else {
		fmt.Printf("Pull Request: %s\n", infosMergeMessage.PRNumber)
		fmt.Printf("Current release version: %s\n", gitVersion)
		fmt.Printf("Patch level: %s\n", patchLevel)
		fmt.Printf("Possible new release version: %s\n", nextVersion)
	}

}
