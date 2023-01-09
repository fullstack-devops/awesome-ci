package scmportal

import (
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"

	log "github.com/sirupsen/logrus"
)

func (lay SCMLayer) GetLatestReleaseVersion() (release *Release, err error) {
	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		latestRelease, err := grc.GetLatestReleaseVersion()
		if err != nil {
			return nil, err
		}
		return &Release{
			TagName:     *latestRelease.TagName,
			Name:        *latestRelease.Name,
			Commit:      *latestRelease.TargetCommitish,
			CreatedAt:   &latestRelease.CreatedAt.Time,
			PublishedAt: &latestRelease.PublishedAt.Time,
		}, nil

	case *gitlab.GitLabRichClient:
		log.Warnln("gitlab is not yet implemented")

		latestRelease, err := grc.GetLatestReleaseVersion()
		if err != nil {
			return nil, err
		}
		return &Release{
			TagName:     latestRelease.TagName,
			Name:        latestRelease.Name,
			Commit:      latestRelease.Commit.String(),
			CreatedAt:   latestRelease.CreatedAt,
			PublishedAt: latestRelease.ReleasedAt,
		}, nil
	}
	return
}
