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
	Long: `Publish a recently created GitHub release

Publish will publish a previously created release to GitHub. The release body
can be written in markdown and will be rendered as such on GitHub. Additionally,
any number of assets can be uploaded, including but not limited to zip and tgz
files. The assets will be uploaded to GitHub and linked to the release.

The assets can be specified as a list of local file paths.

The zip and tgz mechanisms can be used by specifying the path to a directory
containing the files to be uploaded. The directory will be zipped or tarred and
gzipped and uploaded to GitHub as a single asset. The name of the asset will be
the name of the directory with the appropriate extension appended. For example,
if the directory is named "myfiles", the asset will be named "myfiles.zip" or "myfiles.tgz".

Example of using the release publish command

To publish a release with a specific release ID and assets, use the following command:

    awesome-ci release publish --release-id 12345 \
        --asset "file=out/awesome-ci_v1.0.0_amd64" \
        --asset "file=out/awesome-ci_v1.0.0_arm64"
	awesome-ci release publish --release-id 12345 --asset "file=out/awesome-ci_v1.0.0_amd64" --asset "file=out/awesome-ci_v1.0.0_arm64"

If the release ID is not provided, the command will look for the 'ACI_RELEASE_ID' environment variable:

    export ACI_RELEASE_ID=12345
    awesome-ci release publish \
        --asset "file=out/awesome-ci_v1.0.0_amd64" \
        --asset "file=out/awesome-ci_v1.0.0_arm64"
	export ACI_RELEASE_ID=12345
	awesome-ci release publish --asset "file=out/awesome-ci_v1.0.0_amd64" --asset "file=out/awesome-ci_v1.0.0_arm64"

You can also publish a release with a directory as a zip asset:

    awesome-ci release publish --release-id 12345 \
        --asset "zip=out/myfiles"

The assets should be specified as local file paths to be uploaded to the GitHub release.`,
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
			} else {
				log.Fatalln("no release ID provided")
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
	publishCmd.Flags().StringArrayVarP(&assets, "asset", "a", []string{}, "add an asset to the release, can be specified multiple times.")
}
