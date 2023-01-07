package scmportal

import (
	"awesome-ci/internal/app/awesome-ci/ces"
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"awesome-ci/internal/pkg/rcpersist"
)

func LoadSCMPortalLayer() (cesType rcpersist.CESType, grc interface{}, err error) {
	cesType, scmPortalType, connCreds, err := ces.DetectCes()
	if err != nil {
		return
	}

	switch scmPortalType {
	case rcpersist.SCMPortalTypeGitHub:
		ghrc, err := github.NewGitHubClient(&connCreds.ServerUrl, &connCreds.Repository, &connCreds.Token)
		if err != nil {
			return cesType, nil, err
		}
		return cesType, ghrc, nil

	case rcpersist.SCMPortalTypeGitLab:
		glrc, err := gitlab.NewGitLabClient(&connCreds.ServerUrl, &connCreds.Repository, &connCreds.Token)
		if err != nil {
			return cesType, nil, err
		}
		return cesType, glrc, nil
	}

	return
}
