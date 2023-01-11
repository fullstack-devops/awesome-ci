package scmportal

import (
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"

	log "github.com/sirupsen/logrus"
)

// CommentHelpToPullRequest comments to a pull or merge request if awesome ci is running in CI mode
func (lay SCMLayer) CommentHelpToPullRequest(number int) (err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isSilentBool := os.LookupEnv("ACI_SILENT")
	if isCI && !isSilentBool {

		switch grc := lay.Grc.(type) {

		case *github.GitHubRichClient:
			err := grc.CommentHelpToPullRequest(number)
			if err != nil {
				return err
			}

		case *gitlab.GitLabRichClient:
			// not implemented
			log.Warnln("gitlab is not yet implemented")
		}

	}
	return nil
}
