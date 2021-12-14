package gitController

import (
	"github.com/google/go-github/v39/github"
	"github.com/xanzy/go-gitlab"
)

type AciPrInfos struct {
	PrNumber       int
	PatchLevel     string
	CurrentVersion string
	LatestVersion  string
	NextVersion    string
	Sha            string
	ShaShort       string
	BranchName     string
	GithubPrInfos  *github.PullRequest
	GitlabPrInfos  *gitlab.MergeRequest
}
