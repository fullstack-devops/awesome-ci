package acigithub

import (
	"awesome-ci/src/models"
	"awesome-ci/src/semver"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
)

var (
	ctx                                  = context.Background()
	githubRepository, isgithubRepository = os.LookupEnv("GITHUB_REPOSITORY")
)

// GetPrInfos need the PullRequest-Number
func GetPrInfos(prNumber int) (standardPrInfos *models.StandardPrInfos, prInfos *github.PullRequest, err error) {
	if !isgithubRepository {
		log.Fatalln("make shure the GITHUB_REPOSITORY is available!")
	}
	owner, repo := devideOwnerAndRepo(githubRepository)
	if prNumber != 0 {
		prInfos, _, err = GithubClient.PullRequests.Get(ctx, owner, repo, prNumber)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load any information about the given pull request  %d: %v", prNumber, err)
		}
	}

	prSHA := *prInfos.Head.SHA
	branchName := *prInfos.Head.Ref
	patchLevel := branchName[:strings.Index(branchName, "/")]

	// if an comment exists with aci=major, make a major version!
	issueComments, err := GetIssueComments(prNumber, owner, repo)
	if err != nil {
		return nil, nil, err
	}
	for _, comment := range issueComments {
		// Must have permission in the repo to create a major version
		// MANNEQUIN|NONE https://docs.github.com/en/graphql/reference/enums#commentauthorassociation
		if strings.Contains(*comment.Body, "aci=major") && strings.Contains("OWNER|CONTRIBUTOR|COLLABORATOR", *comment.AuthorAssociation) {
			patchLevel = "major"
			break
		}
	}

	repositoryRelease, err := GetLatestReleaseVersion(owner, repo)
	if err != nil {
		return nil, nil, err
	}
	nextVersion, err := semver.IncreaseVersion(patchLevel, *repositoryRelease.TagName)
	if err != nil {
		return nil, nil, err
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
		NextVersion:    nextVersion,
	}
	return
}
