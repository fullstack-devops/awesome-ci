package scmportal

import (
	"os"

	"github.com/layer87-labs/relctl/internal/app/relctl/scm-portal/github"
)

// CommentHelpToPullRequest comments to a pull or merge request if awesome ci is running in CI mode
func (lay SCMLayer) CommentHelpToPullRequest(number int) (err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isSilentBool := os.LookupEnv("RELCTL_SILENT")
	if isCI && !isSilentBool {
		grc := lay.Grc.(*github.GitHubRichClient)
		err := grc.CommentHelpToPullRequest(number)
		if err != nil {
			return err
		}
	}
	return nil
}
