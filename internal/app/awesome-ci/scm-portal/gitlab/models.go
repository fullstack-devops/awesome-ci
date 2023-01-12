package gitlab

import "github.com/xanzy/go-gitlab"

type GitLabRichClient struct {
	Client     *gitlab.Client
	Owner      string
	Repository string
}
