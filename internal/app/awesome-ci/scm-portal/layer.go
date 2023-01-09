package scmportal

import (
	"awesome-ci/internal/app/awesome-ci/ces"
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"awesome-ci/internal/pkg/rcpersist"
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
		ghrc, err := github.NewGitHubClient(&connCreds.ServerUrl, &connCreds.Repository, &connCreds.Token)
		if err != nil {
			return nil, err
		}
		return &SCMLayer{
			CES: ces,
			Grc: ghrc,
		}, err

	case rcpersist.SCMPortalTypeGitLab:
		glrc, err := gitlab.NewGitLabClient(&connCreds.ServerUrl, &connCreds.Repository, &connCreds.Token)
		if err != nil {
			return nil, err
		}
		return &SCMLayer{
			CES: ces,
			Grc: glrc,
		}, err
	}

	return
}
