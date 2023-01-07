package scmportal

import "awesome-ci/internal/pkg/semver"

type PrMrRequestInfos struct {
	Number         int
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
