package service

import (
	"awesome-ci/gitcontroller"
	"awesome-ci/semver"
	"fmt"
	"strings"
)

func GetBuildInfos(cienv string, overrideVersion *string, getVersionIncrease *string, format *string) {
	var gitVersion string
	if *overrideVersion != "" {
		gitVersion = *overrideVersion
	} else {
		gitVersion = gitcontroller.GetLatestReleaseVersion(cienv)
	}

	var infosMergeMessage infosMergeMessage
	if cienv == "Github" {
		var err error
		infosMergeMessage, err = getMergeMessage()
		if err != nil {
			infosMergeMessage.PRNumber = "/"
		}
	}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		patchLevel = infosMergeMessage.PatchLevel
	}

	nextVersion := semver.IncreaseSemVer(patchLevel, gitVersion)
	replacer := strings.NewReplacer(
		"pr", infosMergeMessage.PRNumber,
		"version", nextVersion,
		"patchLevel", patchLevel)
	output := replacer.Replace(*format)
	fmt.Print(output)
}
