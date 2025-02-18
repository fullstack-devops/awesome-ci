package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"
)

func LoadSCMPortalLayer() (scmLayer *SCMLayer, err error) {
	ces, scmPortalType, connCreds, err := ces.DetectCes()
	if err != nil {
		return &SCMLayer{
			CES: ces,
		}, err
	}

	switch scmPortalType {
	case rcpersist.SCMPortalTypeGitHub:
		ghrc, err := github.NewGitHubClient(&connCreds.ServerURL, &connCreds.Repository, &connCreds.Token)
		if err != nil {
			return nil, err
		}
		return &SCMLayer{
			CES: ces,
			Grc: ghrc,
		}, err
	}

	return
}
