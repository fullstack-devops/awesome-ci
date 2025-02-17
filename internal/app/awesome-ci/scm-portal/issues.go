package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"
)

func (lay SCMLayer) SearchIssuesForOverrides(number int) (nextVersion *string, patchLevel *semver.PatchLevel, err error) {
	grc := lay.Grc.(*github.GitHubRichClient)
	nextVersion, patchLevel, err = grc.SearchIssuesForOverrides(number)
	if err != nil {
		return nil, nil, err
	}

	return
}
