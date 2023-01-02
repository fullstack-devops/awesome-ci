package release

import (
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

	},
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a GitHub release",
	Long:  `Print all infos about a pull request in GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	// commands
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(publishCmd)

	// Flags
	Cmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "", false, "make dry-run before writing version to Git by calling it")

	publishCmd.Flags().Int64VarP(&ReleaseId, "releaseid", "", 0, "publish an early defined release")
	publishCmd.Flags().StringVarP(&Assets, "assets", "a", "", "define output by get")
}
