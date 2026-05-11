package scmportal

import (
	"github.com/layer87-labs/relctl/internal/app/relctl/scm-portal/github"
	"github.com/layer87-labs/relctl/internal/pkg/semver"
)

func (lay SCMLayer) SearchIssuesForOverrides(number int) (nextVersion *string, patchLevel *semver.PatchLevel, err error) {
	grc := lay.Grc.(*github.GitHubRichClient)
	nextVersion, patchLevel, err = grc.SearchIssuesForOverrides(number)
	if err != nil {
		return nil, nil, err
	}

	return
}
