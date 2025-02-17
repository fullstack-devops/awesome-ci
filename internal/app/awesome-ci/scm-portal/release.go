package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/uploadasset"

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

func (lay SCMLayer) CreateRelease(tagName string, releasePrefix string, branch string, body string) (createdRelease *Release, err error) {
	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		createdRel, err := grc.CreateRelease(tagName, releasePrefix, branch, body)
		if err != nil {
			return nil, err
		}
		return &Release{
			ID:          *createdRel.ID,
			TagName:     *createdRel.TagName,
			Name:        *createdRel.Name,
			Commit:      *createdRel.TargetCommitish,
			CreatedAt:   &createdRel.CreatedAt.Time,
			PublishedAt: nil,
		}, nil

	case *gitlab.GitLabRichClient:
		log.Fatalln("gitlab is not yet implemented")

		createdRel, err := grc.CreateRelease(tagName, branch, body)
		if err != nil {
			return nil, err
		}
		return &Release{
			ID:          0, // TODO: does not exist at GitLab!
			TagName:     createdRel.TagName,
			Name:        createdRel.Name,
			Commit:      createdRel.Commit.String(),
			CreatedAt:   createdRel.CreatedAt,
			PublishedAt: nil,
		}, nil
	}
	return
}

func (lay SCMLayer) PublishRelease(tagName string, releasePrefix string, branch string, body string, releaseID int64, assets []uploadasset.UploadAsset) (publishedRelease *Release, err error) {
	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		_, err := grc.PublishRelease(tagName, releasePrefix, branch, body, releaseID, assets)
		if err != nil {
			return nil, err
		}
		return nil, nil

	case *gitlab.GitLabRichClient:
		log.Fatalln("gitlab is not yet implemented")

		createdRel, err := grc.CreateRelease(tagName, branch, body)
		if err != nil {
			return nil, err
		}
		return &Release{
			ID:          0, // TODO: does not exist at GitLab!
			TagName:     createdRel.TagName,
			Name:        createdRel.Name,
			Commit:      createdRel.Commit.String(),
			CreatedAt:   createdRel.CreatedAt,
			PublishedAt: createdRel.ReleasedAt,
		}, nil
	}
	return
}
