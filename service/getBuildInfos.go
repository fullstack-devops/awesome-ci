package service

import (
	"awesome-ci/ciRunnerController"
	"awesome-ci/gitOnlineController"
	"awesome-ci/models"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetBuildInfos(cienv string, versionOverr *string, patchLevelOverr *string, format *string) {

	envs, _, buildInfos := setBuildInfos(versionOverr, patchLevelOverr)

	if *format != "" {
		replacer := strings.NewReplacer(
			"pr", fmt.Sprint(buildInfos.PrNumber),
			"version", buildInfos.Version,
			"next_version", buildInfos.NextVersion,
			"patchLevel", buildInfos.PatchLevel)
		output := replacer.Replace(*format)
		fmt.Print(output)
	} else {
		fmt.Println("#### Setting Env variables:")

		for _, env := range envs {
			fmt.Println(env.Key + "=" + env.Value)
		}

		fmt.Println("\n#### Info output:")
		fmt.Printf("Pull Request: %d\n", buildInfos.PrNumber)
		fmt.Printf("Current release version: %s\n", buildInfos.Version)
		fmt.Printf("Patch level: %s\n", buildInfos.PatchLevel)
		fmt.Printf("Possible new release version: %s\n", buildInfos.NextVersion)
	}
}

func setBuildInfos(versionOverr *string, patchLevelOverr *string) (envs []models.BuildEnvironmentVariable, prInfos models.GitHubPullRequest, buildInfos models.BuildInfos) {
	prInfos, prNumber, err := getPRInfos()
	if err != nil {
		panic(err)
	}
	buildInfos.PrNumber = prNumber

	branchName := prInfos.Head.Ref
	buildInfos.PatchLevel = branchName[:strings.Index(branchName, "/")]

	// if an comment exists with aci=major, make a major version!
	if detectIfMajor(buildInfos.PrNumber) {
		buildInfos.PatchLevel = "major"
	}

	if *patchLevelOverr != "" {
		buildInfos.PatchLevel = *patchLevelOverr
	}

	if *versionOverr != "" {
		buildInfos.Version = *versionOverr
	} else {
		buildInfos.Version = gitOnlineController.GetLatestReleaseVersion()
	}
	buildInfos.NextVersion = increaseSemVer(buildInfos.PatchLevel, buildInfos.Version)

	envs = []models.BuildEnvironmentVariable{
		{Key: "ACI_PR", Value: fmt.Sprintf("%d", buildInfos.PrNumber)},
		{Key: "ACI_PR_SHA", Value: prInfos.Head.Sha},
		{Key: "ACI_PR_SHA_SHORT", Value: prInfos.Head.Sha[:8]},
		{Key: "ACI_ORGA", Value: strings.ToLower(CiEnvironment.GitInfos.Orga)},
		{Key: "ACI_REPO", Value: strings.ToLower(CiEnvironment.GitInfos.Repo)},
		{Key: "ACI_BRANCH", Value: branchName},
		{Key: "ACI_PATCH_LEVEL", Value: buildInfos.PatchLevel},
		{Key: "ACI_VERSION", Value: buildInfos.Version},
		{Key: "ACI_NEXT_VERSION", Value: buildInfos.NextVersion},
	}
	ciRunnerController.SetEnvVariables(envs)
	for _, env := range envs {
		os.Setenv(env.Key, env.Value)
	}
	return
}

func getPRInfos() (prInfos models.GitHubPullRequest, prNumber int, err error) {
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
