package scmportal

import (
	"time"

	"github.com/layer87-labs/relctl/internal/app/relctl/ces"
	"github.com/layer87-labs/relctl/internal/pkg/semver"
)

type SCMLayer struct {
	Grc interface{}
	CES ces.CES
}

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

type Release struct {
	ID          int64      `json:"id"`           // GitHub: ID, GitLab: not available!!
	TagName     string     `json:"tag_name"`     // GitHub: PublishedAt; GitLab: ReleasedAt
	Name        string     `json:"name"`         // equaly
	Commit      string     `json:"commit"`       // GitHub: TargetCommitish; GitLab: Commit
	CreatedAt   *time.Time `json:"created_at"`   // equaly
	PublishedAt *time.Time `json:"published_at"` // GitHub: PublishedAt; GitLab: ReleasedAt
}
