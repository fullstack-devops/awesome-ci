package acigithub

import (
	"awesome-ci/src/models"
	"awesome-ci/src/semver"
	"awesome-ci/src/tools"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v39/github"
)

// GetPrInfos need the PullRequest-Number
func GetPrInfos(prNumber int, mergeCommitSha string) (standardPrInfos *models.StandardPrInfos, prInfos *github.PullRequest, err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := tools.DevideOwnerAndRepo(githubRepository)
	if prNumber != 0 {
		prInfos, _, err = GithubClient.PullRequests.Get(ctx, owner, repo, prNumber)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
	}
	if mergeCommitSha != "" && prNumber == 0 {
		prOpts := github.PullRequestListOptions{
			State:     "all",
			Sort:      "updated",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 10,
			},
		}
		pullRequests, _, err := GithubClient.PullRequests.List(ctx, owner, repo, &prOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
		var found int = 0
		for _, pr := range pullRequests {
			if *pr.MergeCommitSHA == mergeCommitSha {
				prInfos = pr
				found = found + 1
			}
		}
		if found > 1 {
			return nil, nil, fmt.Errorf("found more than one pull request, this should not be possible. please open an issue with all log files")
		}
	}

	if prInfos == nil {
		return nil, nil, fmt.Errorf("no pull request found, please check if all resources are specified")
	}

	isCI, isCIBool := os.LookupEnv("CI")
	_, isSilentBool := os.LookupEnv("ACI_SILENT")
	if isCIBool && !isSilentBool {
		if *prInfos.State == "open" && isCI == "true" {
			err = CommentHelpToPullRequest(*prInfos.Number)
			if err != nil {
				log.Println(err)
			}
		}
	}

	prSHA := *prInfos.Head.SHA
	branchName := *prInfos.Head.Ref
	patchLevel := branchName[:strings.Index(branchName, "/")]

	var version = ""
	// if an comment exists with aci_patch_level=major, make a major version!
	issueComments, err := GetIssueComments(prNumber, owner, repo)
	if err != nil {
		return nil, nil, err
	}
	for _, comment := range issueComments {
		// Must have permission in the repo to create a major version
		// MANNEQUIN|NONE https://docs.github.com/en/graphql/reference/enums#commentauthorassociation
		if strings.Contains("OWNER|CONTRIBUTOR|COLLABORATOR", *comment.AuthorAssociation) {
			aciVersionOverride := regexp.MustCompile(`aci_version_override: ([0-9]+\.[0-9]+\.[0-9]+)`)
			aciPatchLevel := regexp.MustCompile(`aci_patch_level: ([a-zA-Z]+)`)

			if aciPatchLevel.MatchString(*comment.Body) {
				patchLevel = aciVersionOverride.FindStringSubmatch(*comment.Body)[1]
				break
			}
			if aciVersionOverride.MatchString(*comment.Body) {
				version = aciVersionOverride.FindStringSubmatch(*comment.Body)[1]
				break
			}
		}
	}

	repositoryRelease, err := GetLatestReleaseVersion(owner, repo)
	if err != nil {
		return nil, nil, err
	}

	if version == "" {
		version, err = semver.IncreaseVersion(patchLevel, *repositoryRelease.TagName)
		if err != nil {
			return nil, nil, err
		}
	}

	standardPrInfos = &models.StandardPrInfos{
		PrNumber:       prNumber,
		Owner:          owner,
		Repo:           repo,
		BranchName:     branchName,
		Sha:            prSHA,
		ShaShort:       prSHA[:8],
		PatchLevel:     patchLevel,
		LatestVersion:  *repositoryRelease.TagName,
		CurrentVersion: "",
		NextVersion:    version,
		MergeCommitSha: *prInfos.MergeCommitSHA,
	}
	return
}
