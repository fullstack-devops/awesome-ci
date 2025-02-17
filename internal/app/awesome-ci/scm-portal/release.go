package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/uploadasset"
)

func (lay SCMLayer) GetLatestReleaseVersion() (release *Release, err error) {

	grc := lay.Grc.(*github.GitHubRichClient)
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
}

func (lay SCMLayer) CreateRelease(tagName string, releasePrefix string, branch string, body string) (createdRelease *Release, err error) {

	grc := lay.Grc.(*github.GitHubRichClient)
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
}

func (lay SCMLayer) PublishRelease(tagName string, releasePrefix string, branch string, body string, releaseID int64, assets []uploadasset.UploadAsset) (publishedRelease *Release, err error) {
	grc := lay.Grc.(*github.GitHubRichClient)
	_, err = grc.PublishRelease(tagName, releasePrefix, branch, body, releaseID, assets)
	if err != nil {
		return nil, err
	}
	return nil, nil

}
