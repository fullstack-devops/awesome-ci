package service

import (
	"awesome-ci/ciRunnerController"
	"awesome-ci/gitOnlineController"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
)

func GetBuildInfos(cienv string, versionOverr *string, patchLevelOverr *string, format *string) {

	prInfos, prNumber, err := getPRInfos()
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
			gitVersion = gitOnlineController.GetLatestReleaseVersion()
		}
	}
	nextVersion := increaseSemVer(patchLevel, gitVersion)

	prSHA := *prInfos.Head.SHA
	var envs []string
	envs = append(envs, fmt.Sprintf("ACI_PR=%d", prNumber))
	envs = append(envs, fmt.Sprintf("ACI_PR_SHA=%s", prSHA))
	envs = append(envs, fmt.Sprintf("ACI_PR_SHA_SHORT=%s", prSHA[:8]))
	envs = append(envs, fmt.Sprintf("ACI_ORGA=%s", strings.ToLower(CiEnvironment.GitInfos.Owner)))
	envs = append(envs, fmt.Sprintf("ACI_REPO=%s", strings.ToLower(CiEnvironment.GitInfos.Repo)))
	envs = append(envs, fmt.Sprintf("ACI_BRANCH=%s", branchName))
	envs = append(envs, fmt.Sprintf("ACI_PATCH_LEVEL=%s", patchLevel))
	envs = append(envs, fmt.Sprintf("ACI_VERSION=%s", gitVersion))
	envs = append(envs, fmt.Sprintf("ACI_NEXT_VERSION=%s", nextVersion))
	ciRunnerController.SetEnvVariables(envs)

	if *format != "" {
		replacer := strings.NewReplacer(
			"pr", fmt.Sprint(prNumber),
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
		fmt.Printf("Pull Request: %d\n", prNumber)
		fmt.Printf("Current release version: %s\n", gitVersion)
		fmt.Printf("Patch level: %s\n", patchLevel)
		fmt.Printf("Possible new release version: %s\n", nextVersion)
	}
}

func getPRInfos() (prInfos *github.PullRequest, prNumber int, err error) {
	prNumber = 0
	// Try to get PR number from 'git name-rev HEAD'
	prNumber, _, err = getNameRevHead()
	if err != nil {
		return
	}
	// Try to get PR number from merge message 'merge'
	if prNumber == 0 {
		prNumber, err = getPrFromMergeMessage()
		if err != nil {
			return
		}
	}
	prInfos, err = gitOnlineController.GetPrInfos(prNumber)
	if err != nil {
		fmt.Println("could not load any information about the current pull request", err)
	}
	return
}

func getPrFromMergeMessage() (pr int, err error) {
	regex := `.*#([0-9]+).*`
	r := regexp.MustCompile(regex)

	mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
	if len(mergeMessage) > 1 {
		return strconv.Atoi(mergeMessage[1])
	} else {
		return 0, errors.New("No PR found in merge message pls make shure this regex matches: " + regex +
			"\nExample: Merge pull request #3 from some-orga/feature/awesome-feature" +
			"\nIf you like to set your patch level manually by flag: -level (feautre|bugfix)")
	}
}

func getDefaultBranch() string {
	return runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
}

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
		pr = gitOnlineController.GetPrNumberForBranch(branchName)
	} else {
		err = errors.New("no branch or pr in 'git name-rev head' found:" + gitNameRevHead)
	}
	return
}

func detectIfMajor(issueNumber int) bool {
	resBool := false
	issueComments, err := gitOnlineController.GetIssueComments(issueNumber)
	if err != nil {
		panic(err)
	}
	for _, comment := range issueComments {
		// Must have permission in the repo to create a major version
		// MANNEQUIN|NONE https://docs.github.com/en/graphql/reference/enums#commentauthorassociation
		if strings.Contains(comment.Body, "aci=major") && !strings.Contains("MANNEQUIN|NONE", comment.AuthorAssociation) {
			resBool = true
			break
		}
	}
	return resBool
}
