package scmportal

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"

	log "github.com/sirupsen/logrus"
)

// GetPrInfos retrieves detailed information about a pull request or merge request
// from the version control system. It uses the provided pull request number or
// merge commit SHA to fetch the relevant information. The function populates a
// PrMrRequestInfos struct with details such as the PR number, SHA, branch name,
// owner, repository, and merge commit SHA. It also parses the patch level from
// the branch name and checks for version or patch level overrides in the issue
// comments. If no version override is found, it determines the latest version
// and calculates the next version based on the patch level.

func (lay SCMLayer) GetPrInfos(number int, mergeCommitSha string) (infos *PrMrRequestInfos, err error) {
	// Initialize the standard information
	infos = &PrMrRequestInfos{}

	// Get the pull request or merge request information
	// from the GitHub or GitLab API

	grc := lay.Grc.(*github.GitHubRichClient)
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

	if version != nil {
		log.Tracef("read comments from pr/mr complete conclusions, found version override version: %s", *version)
	}

	// Check if version override (3)
	if version == nil {
		log.Traceln("no version override found")
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
