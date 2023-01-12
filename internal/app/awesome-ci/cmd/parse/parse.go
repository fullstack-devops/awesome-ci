package parse

import (
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
	Short: "parse json and yaml files and inspect values like jq",
	Long:  `parse json and yaml files and inspect values like jq`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "parse a json string or file",
	Run: func(cmd *cobra.Command, args []string) {
		if File != "" {
			parsejy.ParseFile(Query, File, parsejy.JsonSyntax)
		} else {
			parsejy.Parse(Query, []byte(String), parsejy.JsonSyntax)
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

func init() {
	// commands
	Cmd.AddCommand(jsonCmd)
	Cmd.AddCommand(yamlCmd)

	// Flags
	Cmd.PersistentFlags().StringVarP(&File, "file", "f", "", "file to be parsed")
	Cmd.PersistentFlags().StringVarP(&Query, "query", "q", "", "(required) query for output")
	Cmd.PersistentFlags().StringVarP(&String, "string", "s", "", "query for output")

	Cmd.MarkPersistentFlagRequired("query")
	Cmd.MarkFlagsMutuallyExclusive("file", "string")
	// exclusive Flags
}
