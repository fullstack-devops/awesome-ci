package pullrequest

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/service"

	"github.com/spf13/cobra"
)

var (
	number         int
	formatOut      string
	mergeCommitSha string
)

var Cmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage GitHub pull requests",
	Long:  `The pull request command is used to manage GitHub pull requests. It provides subcommands to get pull request info, allowing you to get all infos about a pull request in GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get pull request info",
	Run: func(cmd *cobra.Command, args []string) {
		service.PrintPRInfos(number, mergeCommitSha, formatOut)
	},
}

func init() {
	// commands
	Cmd.AddCommand(infoCmd)

	// Flags
	infoCmd.Flags().IntVarP(&number, "number", "n", 0, "overwrite the issue number")
	infoCmd.Flags().StringVarP(&mergeCommitSha, "merge-commit-sha", "c", "", "send a given merge commit sha")
	// needs to be implemented
	// infoCmd.Flags().StringVarP(&formatOut, "output", "o", "", "define output available: [json, yaml]")
}
