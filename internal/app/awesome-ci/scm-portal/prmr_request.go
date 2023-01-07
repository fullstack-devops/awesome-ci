package scmportal

import (
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"awesome-ci/internal/pkg/rcpersist"
	"awesome-ci/internal/pkg/semver"

	log "github.com/sirupsen/logrus"
)

func GetPrInfos(cesType rcpersist.CESType, grcInter interface{}, number int, mergeCommitSha string) (infos *PrMrRequestInfos, err error) {
	switch grc := grcInter.(type) {

	case *github.GitHubRichClient:
		prInfos, err := grc.GetPrInfos(number, mergeCommitSha)
		if err != nil {
			return nil, err
		}

		if errCommHelp := CommentHelpToPullRequest(grcInter, number); errCommHelp != nil {
			log.Warnln(errCommHelp)
		}

		prSHA := *prInfos.Head.SHA
		branchName := *prInfos.Head.Ref
		patchLevel := semver.ParsePatchLevel(branchName)

		var nextVersion = ""
		var latestVersion = ""

		infos = &PrMrRequestInfos{
			Number:         number,
			Owner:          grc.Owner,
			Repo:           grc.Repository,
			BranchName:     branchName,
			Sha:            prSHA,
			ShaShort:       prSHA[:8],
			PatchLevel:     patchLevel,
			LatestVersion:  latestVersion,
			NextVersion:    nextVersion,
			MergeCommitSha: *prInfos.MergeCommitSHA,
		}

	case *gitlab.GitLabRichClient:
		// not implemented
	}

	return
}
