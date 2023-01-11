package gitlab

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GetLatestReleaseVersion
func (glrc *GitLabRichClient) GetLatestReleaseVersion() (latestRelease *gitlab.Release, err error) {

	return nil, fmt.Errorf("not implemented")
}

func (glrc *GitLabRichClient) CreateRelease(tagName string, releaseBranch string, body string) (createdRelease *gitlab.Release, err error) {

	return nil, fmt.Errorf("not implemented")
}

func (glrc *GitLabRichClient) PublishRelease(tagName string, releaseBranch string, body string, assets []string) (publishedRelease *gitlab.Release, err error) {

	return nil, fmt.Errorf("not implemented")
}
