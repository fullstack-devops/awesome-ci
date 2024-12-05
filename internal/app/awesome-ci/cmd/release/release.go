package release

import (
	"os"
	"strconv"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/service"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	releaseArgs service.ReleaseArgs
	releaseID   int64
	assets      []string
)

var Cmd = &cobra.Command{
	Use:   "release",
	Short: "Manage GitHub releases with ease",
	Long:  `The release command is used to manage GitHub releases. It provides subcommands to create and publish releases, allowing you to automate the release process and integrate it into CI/CD workflows. Use this command to streamline the release tagging and deployment of your software projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new GitHub release",
	Run: func(cmd *cobra.Command, args []string) {
		service.ReleaseCreate(&releaseArgs)
	},
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a recently created GitHub release",
	Run: func(cmd *cobra.Command, args []string) {

		if releaseID == 0 {
			log.Traceln("looking of env variable 'ACI_RELEASE_ID' since flag is not set")
			releaseIDStr, releaseIDBool := os.LookupEnv("ACI_RELEASE_ID")
			if releaseIDBool {
				releaseIDTmp, err := strconv.ParseInt(releaseIDStr, 10, 64)
				if err != nil {
					log.Fatalln(err)
				} else {
					releaseID = releaseIDTmp
				}
			}
		}

		service.ReleasePublish(&releaseArgs, releaseID, assets)
	},
}

func init() {
	// commands
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(publishCmd)

	// Flags
	Cmd.PersistentFlags().IntVar(&releaseArgs.PrNumber, "prnumber", 0, "overwrite the issue number")
	Cmd.PersistentFlags().BoolVarP(&releaseArgs.DryRun, "dry-run", "", false, "make dry-run before writing version to Git by calling it")
	Cmd.PersistentFlags().BoolVarP(&releaseArgs.Hotfix, "hotfix", "", false, "create a hotfix release")
	Cmd.PersistentFlags().StringVar(&releaseArgs.MergeCommitSHA, "merge-sha", "", "set the merge sha")
	Cmd.PersistentFlags().StringVar(&releaseArgs.ReleaseBranch, "release-branch", "", "set release branch (default: git default)")
	Cmd.PersistentFlags().StringVar(&releaseArgs.ReleasePrefix, "release-prefix", "", "set a custom release prefix (default -> Release or Hotfix)")
	Cmd.PersistentFlags().StringVar(&releaseArgs.Version, "version", "", "override version to Update")
	Cmd.PersistentFlags().StringVarP(&releaseArgs.Body, "body", "b", "", "custom release message (markdown string or file)")
	Cmd.PersistentFlags().StringVarP(&releaseArgs.PatchLevel, "patch-level", "l", "", "predefine patch level of version to Update")

	// exclusive Flags
	publishCmd.Flags().Int64VarP(&releaseID, "release-id", "", 0, "publish an early defined release (also looking for env ACI_RELEASE_ID)")
	publishCmd.Flags().StringArrayVarP(&assets, "asset", "a", []string{}, "define output by get")
}
