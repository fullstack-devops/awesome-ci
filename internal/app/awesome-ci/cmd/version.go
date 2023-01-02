package cmd

import (
	"awesome-ci/internal/app/build"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Awesome-CI",
	Long:  `All software has versions. This is Awesome-CI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", build.Version)
		fmt.Println("Commit: \t", build.CommitHash)
		fmt.Println("Date:   \t", build.BuildDate)
	},
}
