package release

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/service"

	"github.com/spf13/cobra"
)

var (
	Version        string
	PatchLevel     string
	ReleaseId      int64
	Assets         string
	PrNumber       int
	MergeCommitSHA string
	ReleaseBranch  string
	DryRun         bool
	Hotfix         bool
	Body           string
)

var Cmd = &cobra.Command{
	Use:   "release",
	Short: "release",
	Long:  `tbd`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a GitHub release",
	Long:  `Print all infos about a pull request in GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.ReleaseCreate(&service.ReleaseCreateSet{
			Version:        Version,
			PatchLevel:     PatchLevel,
			PrNumber:       PrNumber,
			MergeCommitSHA: MergeCommitSHA,
			ReleaseBranch:  ReleaseBranch,
			DryRun:         DryRun,
			Hotfix:         Hotfix,
			Body:           Body,
		})
	},
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a GitHub release",
	Long:  `Print all infos about a pull request in GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.ReleasePublish(&service.ReleasePublishSet{
			Version:        Version,
			PatchLevel:     PatchLevel,
			PrNumber:       PrNumber,
			MergeCommitSHA: MergeCommitSHA,
			ReleaseBranch:  ReleaseBranch,
			DryRun:         DryRun,
			Hotfix:         Hotfix,
			Body:           Body,
			ReleaseId:      ReleaseId,
			Assets:         Assets,
		})
	},
}

func init() {
	// commands
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(publishCmd)

	// Flags
	Cmd.PersistentFlags().IntVar(&PrNumber, "prnumber", 0, "overwrite the issue number")
	Cmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "", false, "make dry-run before writing version to Git by calling it")
	Cmd.PersistentFlags().BoolVarP(&Hotfix, "hotfix", "", false, "create a hotfix release")
	Cmd.PersistentFlags().StringVar(&MergeCommitSHA, "merge-sha", "", "set the merge sha")
	Cmd.PersistentFlags().StringVar(&ReleaseBranch, "release-branch", "", "set release branch (default: git default)")
	Cmd.PersistentFlags().StringVar(&Version, "version", "", "override version to Update")
	Cmd.PersistentFlags().StringVarP(&Body, "body", "b", "", "custom release message (markdow string or file)")
	Cmd.PersistentFlags().StringVarP(&PatchLevel, "patch-level", "L", "", "predefine patch level of version to Update")

	// exclusive Flags
	publishCmd.Flags().Int64VarP(&ReleaseId, "releaseid", "", 0, "publish an early defined release")
	publishCmd.Flags().StringVarP(&Assets, "assets", "a", "", "define output by get")
}
