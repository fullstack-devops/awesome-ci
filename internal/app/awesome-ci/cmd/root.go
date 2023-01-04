package cmd

import (
	"awesome-ci/internal/app/awesome-ci/cmd/connect"
	"awesome-ci/internal/app/awesome-ci/cmd/parse"
	"awesome-ci/internal/app/awesome-ci/cmd/pullrequest"
	"awesome-ci/internal/app/awesome-ci/cmd/release"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbose bool

var RootCmd = &cobra.Command{
	Use:   "awesome-ci",
	Short: "Awesome CI make your release tagging easy",
	Long: `Awesome CI make your release tagging easy
      Comatible with CI pipelines like Jenkins and GitHub
      Find more information and examples at: https://github.com/fullstack-devops/awesome-ci`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// commands
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(pullrequest.Cmd)
	RootCmd.AddCommand(release.Cmd)
	RootCmd.AddCommand(parse.Cmd)
	RootCmd.AddCommand(connect.Cmd)

	// flags
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// PreRuns
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if Verbose {
			log.SetLevel(log.TraceLevel)
		}
		return nil
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
