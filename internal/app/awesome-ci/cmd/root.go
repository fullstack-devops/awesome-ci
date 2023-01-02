package cmd

import (
	"awesome-ci/internal/app/awesome-ci/cmd/pullrequest"
	"awesome-ci/internal/app/awesome-ci/cmd/release"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "awesome-ci",
	Short: "Awesome CI make your release tagging easy",
	Long: `Awesome CI make your release tagging easy
                Comatible with CI pipelines like Jenkins and GitHub
                Find more information and examples at: https://github.com/fullstack-devops/awesome-ci`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	// commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(pullrequest.Cmd)
	rootCmd.AddCommand(release.Cmd)

	// flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
