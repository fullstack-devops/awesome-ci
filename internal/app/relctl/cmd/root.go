package cmd

import (
	"fmt"
	"os"

	"github.com/layer87-labs/relctl/internal/app/build"
	"github.com/layer87-labs/relctl/internal/app/relctl/cmd/connect"
	"github.com/layer87-labs/relctl/internal/app/relctl/cmd/parse"
	"github.com/layer87-labs/relctl/internal/app/relctl/cmd/pullrequest"
	"github.com/layer87-labs/relctl/internal/app/relctl/cmd/release"
	"github.com/layer87-labs/relctl/internal/app/relctl/cmd/transform"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbose bool

var RootCmd = &cobra.Command{
	Use:     "relctl",
	Version: fmt.Sprintf("%s (%s)-(%s)", build.Version, build.CommitHash, build.BuildDate),
	Short:   "relctl make your release tagging easy",
	Long: `relctl make your release tagging easy
Comatible with CI pipelines like Jenkins and GitHub
Find more information and examples at: https://github.com/layer87-labs/relctl`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	// commands
	RootCmd.AddCommand(pullrequest.Cmd)
	RootCmd.AddCommand(release.Cmd)
	RootCmd.AddCommand(parse.Cmd)
	RootCmd.AddCommand(connect.Cmd)
	RootCmd.AddCommand(transform.Cmd)

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
