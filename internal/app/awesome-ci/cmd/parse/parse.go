package parse

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/parse"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/parsejy"

	"github.com/spf13/cobra"
)

var (
	File   string
	Query  string
	String string
)

var Cmd = &cobra.Command{
	Use:   "parse",
	Short: "parse json and yaml files and inspect values like jq + more",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "parse a json string or file",
	Run: func(cmd *cobra.Command, args []string) {
		if File != "" {
			parsejy.ParseFile(Query, File, parsejy.JSONSyntax)
		} else {
			parsejy.Parse(Query, []byte(String), parsejy.JSONSyntax)
		}
	},
}

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "parse a yaml string or file",
	Run: func(cmd *cobra.Command, args []string) {
		if File != "" {
			parsejy.ParseFile(Query, File, parsejy.YamlSyntax)
		} else {
			parsejy.Parse(Query, []byte(String), parsejy.YamlSyntax)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Parse the given version string and validate against semver",
	Run: func(cmd *cobra.Command, args []string) {
		if err := parse.ValidateVersion(args[0]); err != nil {
			panic(err)
		}
	},
}

func init() {
	// commands
	Cmd.AddCommand(jsonCmd)
	Cmd.AddCommand(yamlCmd)
	Cmd.AddCommand(versionCmd)

	// Flags
	jsonCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "file to be parsed")
	jsonCmd.PersistentFlags().StringVarP(&Query, "query", "q", "", "(required) query for output")
	jsonCmd.PersistentFlags().StringVarP(&String, "string", "s", "", "query for output")

	yamlCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "file to be parsed")
	yamlCmd.PersistentFlags().StringVarP(&Query, "query", "q", "", "(required) query for output")
	yamlCmd.PersistentFlags().StringVarP(&String, "string", "s", "", "query for output")

	jsonCmd.MarkPersistentFlagRequired("query")
	jsonCmd.MarkFlagsMutuallyExclusive("file", "string")

	yamlCmd.MarkPersistentFlagRequired("query")
	yamlCmd.MarkFlagsMutuallyExclusive("file", "string")
	// exclusive Flags
}
