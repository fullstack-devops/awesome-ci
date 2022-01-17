package service

import (
	"errors"
	"regexp"
	"strconv"
)

/* func GetBuildInfos(cienv string, versionOverr *string, patchLevelOverr *string, format *string) {

	prInfos, prNumber, err := getPRInfos(nil)
	if err != nil {
		panic(err)
	}

	branchName := *prInfos.Head.Ref
	patchLevel := branchName[:strings.Index(branchName, "/")]

	// if an comment exists with aci=major, make a major version!
	if detectIfMajor(prNumber) {
		patchLevel = "major"
	}

	if *patchLevelOverr != "" {
		patchLevel = *patchLevelOverr
	}

	var gitVersion string
	if strings.Contains(*format, "version") || *format == "" {
		if *versionOverr != "" {
			gitVersion = *versionOverr
		} else {
			repositoryRelease, err := CiEnvironment.GetLatestReleaseVersion()
			if err != nil {
				log.Println(err)
			}
			gitVersion = *repositoryRelease.TagName
		}
	}
	nextVersion, err := semver.IncreaseVersion(patchLevel, gitVersion)
	.....
} */

func getNameRevHead() (pr int, branchName string, err error) {
	pr = 0
	branchName = ""
	gitNameRevHead := runcmd(`git name-rev HEAD`, true)

	regexIsPR := regexp.MustCompile(`HEAD remotes/pull/([0-9]+)/.*`)
	regexIsBranch := regexp.MustCompile(`HEAD (.*)`)

	regexIsPRMached := regexIsPR.FindStringSubmatch(gitNameRevHead)
	regexIsBranchMached := regexIsBranch.FindStringSubmatch(gitNameRevHead)
	if len(regexIsPRMached) > 1 {
		pr, err = strconv.Atoi(regexIsPRMached[1])
	} else if len(regexIsBranchMached) > 1 {
		branchName = regexIsBranchMached[1]
		// pr, err = CiEnvironment.GetPrNumberForBranch(branchName)
	} else {
		err = errors.New("no branch or pr in 'git name-rev head' found:" + gitNameRevHead)
	}
	return
}
