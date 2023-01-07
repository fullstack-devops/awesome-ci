package gitlab

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func devideOwnerAndRepo(fullRepo string) (owner string, repo string) {
	return strings.Split(fullRepo, "/")[0], strings.Split(fullRepo, "/")[1]
}

// GetPrInfos need the PullRequest-Number
func (glrc *GitLabRichClient) GetMrInfos(mrNumber int) (prInfos *gitlab.MergeRequest, err error) {
	if mrNumber != 0 {
		prInfos, _, err = glrc.Client.MergeRequests.GetMergeRequest(1, mrNumber, nil, nil)
		if err != nil {
			err = errors.New(fmt.Sprintln("could not load any information about the current pull request", err))
			return
		}
	}
	return
}
