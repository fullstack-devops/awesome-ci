package models

import (
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
)

type CIEnvironment struct {
	GitType  string
	GitInfos struct {
		ApiUrl            string
		ApiToken          string
		Repo              string
		Owner             string
		IsOrg             bool
		FullRepo          string
		DefaultBranchName string
	}
	GithubClient *github.Client
	GitlabClient *gitlab.Client
	RunnerType   string
	RunnerInfo   struct {
		EnvFile string
	}
}
