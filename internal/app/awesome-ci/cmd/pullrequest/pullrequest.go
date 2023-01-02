package pullrequest

import (
	"awesome-ci/internal/app/awesome-ci/service"

	"github.com/spf13/cobra"
)

var (
	number int
	format string
)

var Cmd = &cobra.Command{
	Use:   "pr",
	Short: "pull request",
	Long:  `All software has versions. This is Awesome-CI's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "get pull request info",
	Long:  `Print all infos about a pull request in GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.PrintPRInfos(&service.PullRequestInfoSet{
			Number: number,
			Format: format,
		})
	},
}

func init() {
	// commands
	Cmd.AddCommand(infoCmd)

	// Flags
	infoCmd.Flags().IntVarP(&number, "number", "n", 0, "overwrite the issue number")
	infoCmd.Flags().StringVarP(&format, "format", "f", "", "define output by get")
}
