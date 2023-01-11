package connect

import (
	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/connect"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"

	"github.com/spf13/cobra"
)

var (
	Token      string
	Server     string
	Repository string
)

var Cmd = &cobra.Command{
	Use:   "connect",
	Short: "connect persistent to github by creating an encrypted rc file",
	Long: `connect initial to a GitHub or GitHub Enterprise
				at some point you can also connect to GitLab (not yet implemented)
				this is only useful without a runner or in an jenkins pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "connect initial to a GitHub or GitHub Enterprise",
	Run: func(cmd *cobra.Command, args []string) {
		connect.UpdateRcFileForGitHub(Server, Repository, Token)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove all connection files and secrets",
	Run: func(cmd *cobra.Command, args []string) {
		rcpersist.RemoveRcFile()
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check current connection",
	Run: func(cmd *cobra.Command, args []string) {
		connect.CheckConnection()
	},
}

func init() {
	// commands
	Cmd.AddCommand(githubCmd)
	Cmd.AddCommand(removeCmd)
	Cmd.AddCommand(checkCmd)

	// Flags
	githubCmd.Flags().StringVarP(&Server, "server", "s", "https://github.com", "(required) github instance to connect")
	githubCmd.Flags().StringVarP(&Repository, "repository", "r", "", "(required) repo eg.: octo-org/octo-repo")
	githubCmd.Flags().StringVarP(&Token, "token", "t", "", "(required) plain token eg.: ghp_*****")

	githubCmd.MarkFlagRequired("repository")
	githubCmd.MarkFlagRequired("token")

	// exclusive Flags
}
