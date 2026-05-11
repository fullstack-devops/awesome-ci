package parse

import (
	"github.com/layer87-labs/relctl/internal/app/relctl/parse"
	"github.com/layer87-labs/relctl/internal/pkg/parsejy"

	"github.com/spf13/cobra"
)

var (
	File   string
	Query  string
	String string
)

var Cmd = &cobra.Command{
	Use:   "parse",
	Short: "inspect and parse JSON and YAML files to retrieve values, similar to jq with additional features",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "parse a json string or file",
	Run: func(cmd *cobra.Command, args []string) {
		if File != "" {
			_ = parsejy.ParseFile(Query, File, parsejy.JSONSyntax)
		} else {
			_ = parsejy.Parse(Query, []byte(String), parsejy.JSONSyntax)
		}
	},
}

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "parse a yaml string or file",
	Run: func(cmd *cobra.Command, args []string) {
		if File != "" {
			_ = parsejy.ParseFile(Query, File, parsejy.YamlSyntax)
		} else {
			_ = parsejy.Parse(Query, []byte(String), parsejy.YamlSyntax)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Validate the given version string against semver syntax",
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
	jsonCmd.Flags().StringVarP(&File, "file", "f", "", "file to be parsed")
	jsonCmd.Flags().StringVarP(&Query, "query", "q", "", "(required) query for output")
	jsonCmd.Flags().StringVarP(&String, "string", "s", "", "query for output")

	yamlCmd.Flags().StringVarP(&File, "file", "f", "", "file to be parsed")
	yamlCmd.Flags().StringVarP(&Query, "query", "q", "", "(required) query for output")
	yamlCmd.Flags().StringVarP(&String, "string", "s", "", "query for output")

	_ = jsonCmd.MarkFlagRequired("query")
	jsonCmd.MarkFlagsMutuallyExclusive("file", "string")
	_ = jsonCmd.MarkFlagFilename("file")

	_ = yamlCmd.MarkFlagRequired("query")
	yamlCmd.MarkFlagsMutuallyExclusive("file", "string")
	_ = yamlCmd.MarkFlagFilename("file")
	// exclusive Flags
}
