package transform

import (
	"fmt"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/parsejy"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	GroupPrefixInt  int
	SubstractString int
)

var Cmd = &cobra.Command{
	Use:     "transform",
	Aliases: []string{"tf"},
	Short:   "transform given input to json",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var groupByCmd = &cobra.Command{
	Use:     "group-by",
	Short:   "group string array by prefix to json",
	Example: "awesome-ci tf group-by -p 3 --sub-prefix 4 st1_infrastructure-base  st2_ubi9-openjdk-11  st2_ubi9-openjdk-17 \n   produces: {\"st1\":[\"infrastructure-base\"],\"st2\":[\"ubi9-openjdk-11\",\"ubi9-openjdk-17\"]}",
	Run: func(cmd *cobra.Command, args []string) {
		if res, err := parsejy.GroupByPrefix(args, GroupPrefixInt, SubstractString); err != nil {
			logrus.Fatalln(err)
		} else {
			fmt.Println(res)
		}
	},
}

func init() {
	// commands
	Cmd.AddCommand(groupByCmd)

	// Flags
	groupByCmd.PersistentFlags().IntVarP(&GroupPrefixInt, "prefix", "p", 3, "group by prefix until index -- eg.: 'st1_base-image' and int 3, will be grouped by 'st1'")
	groupByCmd.PersistentFlags().IntVarP(&SubstractString, "sub-prefix", "", 0, "remove prefix until index number -- eg.: 'st1_base-image' and 4 will be 'base-image'")

	groupByCmd.MarkPersistentFlagRequired("prefix")
	// exclusive Flags
}
