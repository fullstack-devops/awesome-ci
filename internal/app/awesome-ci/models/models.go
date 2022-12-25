package models

import "awesome-ci/internal/app/awesome-ci/semver"

type StandardPrInfos struct {
	PrNumber       int
	Owner          string
	Repo           string
	PatchLevel     semver.PatchLevel
	CurrentVersion string
	LatestVersion  string
	NextVersion    string
	Sha            string
	ShaShort       string
	BranchName     string
	MergeCommitSha string
}