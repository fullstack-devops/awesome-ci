package scmportal

import (
	"awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"awesome-ci/internal/app/awesome-ci/scm-portal/gitlab"
	"fmt"
	"os"
)

// CommentHelpToPullRequest comments to a pull or merge request if awesome ci is running in CI mode
func CommentHelpToPullRequest(grcInter interface{}, number int) (err error) {
	ci, ciPresent := os.LookupEnv("CI")
	isCI := ci == "true" && ciPresent

	_, isSilentBool := os.LookupEnv("ACI_SILENT")
	if isCI && !isSilentBool {

		switch grc := grcInter.(type) {

		case *github.GitHubRichClient:
			fmt.Println("Debug!!! make help comment")
			err := grc.CommentHelpToPullRequest(number)
			if err != nil {
				return err
			}

		case *gitlab.GitLabRichClient:
			// not implemented
		}

	}
	return nil
}
