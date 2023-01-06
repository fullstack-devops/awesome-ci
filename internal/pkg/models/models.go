package models

import "awesome-ci/internal/pkg/semver"

type StandardPrInfos struct {
	PrNumber       int
	Owner          string
	Repo           string
	PatchLevel     semver.PatchLevel
	LatestVersion  string
	NextVersion    string
	Sha            string
	ShaShort       string
	BranchName     string
	MergeCommitSha string
}

type StandardConnectCredentials struct {
	ServerUrl  string
	Repository string
	Token      string
}
