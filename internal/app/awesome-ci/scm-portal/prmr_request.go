package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"

	log "github.com/sirupsen/logrus"
)

// GetPrInfos retrieves pull request or merge request information
// from the GitHub or GitLab API, comments help instructions to the
// issue and reads comments from the pull request or merge request
// to look for overrides.
//
// The function returns the standard information about the pull
// request or merge request and the next version to be released
// based on the overrides found in the comments.
//
// If no version override is found, the function retrieves the
// latest release version from the GitHub API and increments the
// version based on the patch level override.
//
// The function returns the standard information and the next
// version to be released.
func (lay SCMLayer) GetPrInfos(number int, mergeCommitSha string) (infos *PrMrRequestInfos, err error) {
	// Initialize the standard information
	infos = &PrMrRequestInfos{}

	// Get the pull request or merge request information
	// from the GitHub or GitLab API
	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		prInfos, err := grc.GetPrInfos(number, mergeCommitSha)
		if err != nil {
			return nil, err
		}
		log.Traceln("got PR information from scm-portal-layer")

		infos.Number = *prInfos.Number
		infos.Sha = *prInfos.Head.SHA
		infos.ShaShort = infos.Sha[:8]
		infos.BranchName = *prInfos.Head.Ref
		infos.Owner = grc.Owner
		infos.Repo = grc.Repository
		infos.MergeCommitSha = *prInfos.MergeCommitSHA

		// Parse the patch level from the branch name
		if infos.PatchLevel, err = semver.ParsePatchLevelFormBranch(infos.BranchName); err != nil {
			log.Warnln(err)
		}
		log.Traceln("completed PR standard information to: ", *infos)

	case *gitlab.GitLabRichClient:
		// GitLab is not yet implemented
		log.Warnln("gitlab is not yet implemented")
	}

	// Comment help instructions to the issue
	log.Traceln("comment help instructions to issue")
	if errCommHelp := lay.CommentHelpToPullRequest(infos.Number); errCommHelp != nil {
		log.Warnln(errCommHelp)
	}
	log.Traceln("commented help instructions to issue")

	// Read comments from the pull request or merge request
	// and look for overrides
	log.Traceln("read comments from pr/mr and looking for overrides")
	version, patchLevel, err := lay.SearchIssuesForOverrides(infos.Number)
	if err != nil {
		return nil, err
	}
	log.Tracef("read comments from pr/mr complete conclusions (if nil no override), version: %s, patchLevel: %v", *version, patchLevel)

	// Check if version override (3)
	if version == nil {
		// Get the latest release, if any
		log.Traceln("get latest release, if any")
		repositoryRelease, err := lay.GetLatestReleaseVersion()
		if err != nil {
			log.Infoln("no github release found -> writing default 0.0.0")
			infos.LatestVersion = "0.0.0"
		} else {
			log.Infoln("found latest release", repositoryRelease.TagName)
			infos.LatestVersion = repositoryRelease.TagName
		}

		if patchLevel != nil {
			infos.PatchLevel = *patchLevel
			log.Infof("detected a patch level override to %s", infos.PatchLevel)
		}
		// Increment the version based on the patch level override
		if infos.NextVersion, err = semver.IncreaseVersion(infos.PatchLevel, infos.LatestVersion); err != nil {
			return nil, err
		}

		return infos, nil

	} else {
		log.Infoln("version override via pr comments specified")
		infos.NextVersion = *version
	}

	return
}
