package service

import (
	"awesome-ci/ciRunnerController"
	"awesome-ci/gitOnlineController"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func GetBuildInfos(cienv string, overrideVersion *string, getVersionIncrease *string, format *string) {

	var infosMergeMessage infosMergeMessage
	pr, branchName, err := getNameRevHead()
	if err != nil {
		panic(err)
	}

	if pr != 0 {
		prInfos, err := gitOnlineController.GetPrInfos(pr)
		if err != nil {
			fmt.Println("could not load any information about the current pull request", err)
		}
		branchName = prInfos.Head.Ref
	}
	//if cienv == "Github" {
	//	var err error
	infosMergeMessage, err = getLatestCommitMessage()
	if err != nil {

	}
	//}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		patchLevel = infosMergeMessage.PatchLevel
		if patchLevel == "" {
			i := strings.Index(branchName, "/")
			patchLevel = branchName[:i]
		}
	}

	var gitVersion string
	if strings.Contains(*format, "version") || *format == "" {
		if *overrideVersion != "" {
			gitVersion = *overrideVersion
		} else {
			gitVersion = gitOnlineController.GetLatestReleaseVersion()
		}
	}
	nextVersion := increaseSemVer(patchLevel, gitVersion)

	var envs []string
	envs = append(envs, fmt.Sprintf("PR=%s", infosMergeMessage.PRNumber))
	envs = append(envs, fmt.Sprintf("ORGA=%s", CiEnvironment.GitInfos.Orga))
	envs = append(envs, fmt.Sprintf("REPO=%s", CiEnvironment.GitInfos.Repo))
	envs = append(envs, fmt.Sprintf("VERSION=%s", gitVersion))
	envs = append(envs, fmt.Sprintf("NEXT_VERSION=%s", nextVersion))
	ciRunnerController.SetEnvVariables(envs)

	if *format != "" {
		replacer := strings.NewReplacer(
			"pr", infosMergeMessage.PRNumber,
			"version", gitVersion,
			"next_version", nextVersion,
			"patchLevel", patchLevel)
		output := replacer.Replace(*format)
		fmt.Print(output)
	} else {
		fmt.Println("#### Setting Env variables:")

		for _, env := range envs {
			fmt.Println(env)
		}

		fmt.Println("\n#### Info output:")
		fmt.Printf("Pull Request: %s\n", infosMergeMessage.PRNumber)
		fmt.Printf("Current release version: %s\n", gitVersion)
		fmt.Printf("Patch level: %s\n", patchLevel)
		fmt.Printf("Possible new release version: %s\n", nextVersion)
	}
}

func getLatestCommitMessage() (infos infosMergeMessage, err error) {
	// Output: []string {FullString, PR, FullBranch, Orga, branch, branchBegin, restOfBranch}
	regex := `[a-zA-z ]+#([0-9]+) from (([0-9a-zA-Z\-]+)/(([0-9a-z\-]+)/(.+)))`
	r := regexp.MustCompile(regex)

	// mergeMessage := r.FindStringSubmatch(`Merge pull request #3 from test-orga/feature/test-1`)
	mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
	if len(mergeMessage) > 0 {
		infos.PRNumber = mergeMessage[1]
		infos.PatchLevel = mergeMessage[5]
		return infos, nil
	} else {
		return infos, errors.New("No merge message found pls make shure this regex matches: " + regex +
			"\nExample: Merge pull request #3 from some-orga/feature/awesome-feature" +
			"\nIf you like to set your patch level manually by flag: -level (feautre|bugfix)")
	}
}

func getDefaultBranch() string {
	return runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
}

func getCurrentBranchName() string {
	branchName := runcmd(`git name-rev HEAD`, true)
	return strings.ReplaceAll(branchName, "\n", "")
}

func getNameRevHead() (pr int, branchName string, err error) {
	pr = 0
	branchName = ""
	gitNameRevHead := runcmd(`git name-rev HEAD`, true)

	regexIsPR := regexp.MustCompile(`HEAD remotes/pull/([0-9]+)/.*`)
	regexIsBranch := regexp.MustCompile(`HEAD (.*)`)

	regexIsPRMached := regexIsPR.FindStringSubmatch(gitNameRevHead)
	regexIsBranchMached := regexIsBranch.FindStringSubmatch(gitNameRevHead)
	if len(regexIsPRMached) > 2 {
		pr, err = strconv.Atoi(regexIsPRMached[1])
	} else if len(regexIsBranchMached) > 2 {
		branchName = regexIsPRMached[1]
	} else {
		err = errors.New("no branch or pr in 'git name-rev head' found")
	}
	return
}
