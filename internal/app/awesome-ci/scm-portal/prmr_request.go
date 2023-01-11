package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"

	log "github.com/sirupsen/logrus"
)

// GetPrInfos
// 1. get pr/mr infos from github or gitlab
// 2. comment help instructions to issue
// 3. read comments from pr/mr and looking for overrides
func (lay SCMLayer) GetPrInfos(number int, mergeCommitSha string) (infos *PrMrRequestInfos, err error) {
	infos = &PrMrRequestInfos{}
	infos.Number = number

	// 1. get pr/mr infos from github or gitlab
	switch grc := lay.Grc.(type) {

	case *github.GitHubRichClient:
		prInfos, err := grc.GetPrInfos(number, mergeCommitSha)
		if err != nil {
			return nil, err
		}

		infos.Sha = *prInfos.Head.SHA
		infos.ShaShort = infos.Sha[:8]
		infos.BranchName = *prInfos.Head.Ref
		infos.Owner = grc.Owner
		infos.Repo = grc.Repository
		infos.MergeCommitSha = *prInfos.MergeCommitSHA

		if infos.PatchLevel, err = semver.ParsePatchLevelFormBranch(infos.BranchName); err != nil {
			log.Warnln(err)
		}

	case *gitlab.GitLabRichClient:
		// not implemented
		log.Warnln("gitlab is not yet implemented")
	}

	// 2. comment help instructions to issue
	if errCommHelp := lay.CommentHelpToPullRequest(number); errCommHelp != nil {
		log.Warnln(errCommHelp)
	}

	// 3. read comments from pr/mr and looking for overrides
	version, patchLevel, err := lay.SearchIssuesForOverrides(number)
	if err != nil {
		return nil, err
	}

	// get latest release, if any
	repositoryRelease, err := lay.GetLatestReleaseVersion()
	if err != nil {
		log.Infoln("no github release found -> wirting default 0.0.0")
		infos.LatestVersion = "0.0.0"
	} else {
		infos.LatestVersion = repositoryRelease.TagName
	}

	// check if version override (3)
	if version == nil {

		if patchLevel != nil {
			infos.PatchLevel = *patchLevel
			log.Infof("detected a patch level override to %s", infos.PatchLevel)
		}
		if infos.NextVersion, err = semver.IncreaseVersion(infos.PatchLevel, infos.LatestVersion); err != nil {
			return nil, err
		}

		return infos, nil

	} else {
		log.Traceln("version override via pr comments specified")
		infos.NextVersion = *version
	}

	return
}
