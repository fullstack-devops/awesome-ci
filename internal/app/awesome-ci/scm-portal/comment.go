package scmportal

import (
	"os"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
)

// CommentHelpToPullRequest comments to a pull or merge request if awesome ci is running in CI mode
func (lay SCMLayer) CommentHelpToPullRequest(number int) (err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isSilentBool := os.LookupEnv("ACI_SILENT")
	if isCI && !isSilentBool {
		grc := lay.Grc.(*github.GitHubRichClient)
		err := grc.CommentHelpToPullRequest(number)
		if err != nil {
			return err
		}
	}
	return nil
}
