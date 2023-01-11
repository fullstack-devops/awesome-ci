package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"

	log "github.com/sirupsen/logrus"
)

func (lay SCMLayer) SearchIssuesForOverrides(number int) (nextVersion *string, patchLevel *semver.PatchLevel, err error) {
	var nVersion string
	var pLevel semver.PatchLevel

	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		nVersion, pLevel, err = grc.SearchIssuesForOverrides(number)
		if err != nil {
			return nil, nil, err
		}

	case *gitlab.GitLabRichClient:
		// not implemented
		log.Warnln("gitlab is not yet implemented")
	}

	if nVersion == "" {
		nextVersion = nil
	} else {
		nextVersion = &nVersion
	}
	if pLevel == "" {
		patchLevel = nil
	} else {
		patchLevel = &pLevel
	}

	return
}
