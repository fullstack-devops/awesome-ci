package gitlabapi

import (
	"awesome-ci/internal/pkg/models"
	"awesome-ci/internal/pkg/semver"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

var (
	ctx context.Context
	// githubRepository, isgithubRepository = os.LookupEnv("GITHUB_REPOSITORY")
)

func devideOwnerAndRepo(fullRepo string) (owner string, repo string) {
	return strings.Split(fullRepo, "/")[0], strings.Split(fullRepo, "/")[1]
}

// GetPrInfos need the PullRequest-Number
func GetMrInfos(mrNumber int) (standardPrInfos *models.StandardPrInfos, prInfos *gitlab.MergeRequest, err error) {
	if mrNumber != 0 {
		prInfos, _, err = GitLabClient.MergeRequests.GetMergeRequest(1, mrNumber, nil, nil)
		if err != nil {
			err = errors.New(fmt.Sprintln("could not load any information about the current pull request", err))
			return
		}
	}

	prSHA := prInfos.SHA
	branchName := prInfos.Reference
	patchLevel := semver.ParsePatchLevel(branchName)

	standardPrInfos = &models.StandardPrInfos{
		PrNumber:   mrNumber,
		BranchName: branchName,
		Sha:        prSHA,
		ShaShort:   prSHA[:8],
		PatchLevel: patchLevel,
	}
	return
}
