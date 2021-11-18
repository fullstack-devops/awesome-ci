package service

import (
	"awesome-ci/src/gitController"
	"awesome-ci/src/semver"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v39/github"
)

type PullRequestSet struct {
	Fs   *flag.FlagSet
	Info PullRequestInfoSet
}

type PullRequestInfoSet struct {
	Fs         *flag.FlagSet
	Number     int
	EvalNumber bool
}

func PrintPRInfos(args *PullRequestInfoSet) {
	prInfos, _, err := getPRInfos(args)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(prInfos)
}

func getPRInfos(args *PullRequestInfoSet) (prInfos *github.PullRequest, prNumber int, err error) {
	if args.Number != 0 {
		prInfos, err = CiEnvironment.GetPrInfos(args.Number)
		if err != nil {
			fmt.Println("could not load any information about the current pull request", err)
		}
	}

	branchName := *prInfos.Head.Ref
	patchLevel := branchName[:strings.Index(branchName, "/")]

	// if an comment exists with aci=major, make a major version!
	if detectIfMajor(prNumber) {
		patchLevel = "major"
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

	return
}

func setPrInfosToEnv(prInfos *github.PullRequest) {
	prSHA := *prInfos.Head.SHA
	var envs []string
	envs = append(envs, fmt.Sprintf("ACI_PR=%d", *prInfos.Number))
	envs = append(envs, fmt.Sprintf("ACI_PR_SHA=%s", prSHA))
	envs = append(envs, fmt.Sprintf("ACI_PR_SHA_SHORT=%s", prSHA[:8]))
	envs = append(envs, fmt.Sprintf("ACI_ORGA=%s", strings.ToLower(*CiEnvironment.GitInfos.Owner)))
	envs = append(envs, fmt.Sprintf("ACI_REPO=%s", strings.ToLower(*CiEnvironment.GitInfos.Repo)))
	envs = append(envs, fmt.Sprintf("ACI_BRANCH=%s", *prInfos.Head.Ref))
	envs = append(envs, fmt.Sprintf("ACI_PATCH_LEVEL=%s", patchLevel))
	envs = append(envs, fmt.Sprintf("ACI_VERSION=%s", gitVersion))
	envs = append(envs, fmt.Sprintf("ACI_NEXT_VERSION=%s", nextVersion))
	gitController.SetEnvVariables(envs)
}

/* func getPRInfos() (prInfos *github.PullRequest, prNumber int, err error) {
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
	prInfos, err = CiEnvironment.GetPrInfos(prNumber)
	if err != nil {
		fmt.Println("could not load any information about the current pull request", err)
	}
	return
} */
