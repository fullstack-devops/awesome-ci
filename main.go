package main

import (
	"awesome-ci/src/service"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	// cienv       string
	releaseSet     service.ReleaseSet
	pullrequestSet service.PullRequestSet
	parseSet       ParseSet
	version        string
	versionFlag    bool
)

type ParseSet struct {
	Fs   *flag.FlagSet
	Json struct {
		Fs    *flag.FlagSet
		file  string
		value string
	}
	Yaml struct {
		Fs    *flag.FlagSet
		file  string
		value string
	}
}

func init() {
	flag.BoolVar(&versionFlag, "version", false, "print version by calling it")

	// PullRequestSet
	pullrequestSet.Fs = flag.NewFlagSet("pr", flag.ExitOnError)
	pullrequestSet.Fs.Usage = func() {
		fmt.Println("Available commands:")
		fmt.Println("  info")
		fmt.Println("Use \"awesome-ci pr <command> --help\" for more information about a given command.")
	}
	pullrequestSet.Info.Fs = flag.NewFlagSet("pr info", flag.ExitOnError)
	pullrequestSet.Info.Fs.IntVar(&pullrequestSet.Info.Number, "number", 0, "overwrite the issue number")
	pullrequestSet.Info.Fs.StringVar(&pullrequestSet.Info.Format, "format", "", "define output by get")

	// ReleaseSet
	releaseSet.Fs = flag.NewFlagSet("release", flag.ExitOnError)
	releaseSet.Fs.Usage = func() {
		fmt.Println("Available commands:")
		// fmt.Println("  info")
		fmt.Println("  create")
		fmt.Println("  publish")
		fmt.Println("Use \"awesome-ci release <command> --help\" for more information about a given command.")
	}
	releaseSet.Create.Fs = flag.NewFlagSet("release create", flag.ExitOnError)
	releaseSet.Create.Fs.IntVar(&releaseSet.Create.PrNumber, "prnumber", 0, "overwrite the issue number")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.MergeCommitSHA, "merge-sha", "", "set the merge sha")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.ReleaseBranch, "release-branch", "", "set release branch (default: git default)")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.Version, "version", "", "override version to Update")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.Body, "body", "", "custom release message (markdow string or file)")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.PatchLevel, "patchLevel", "", "predefine version to Update")
	releaseSet.Create.Fs.BoolVar(&releaseSet.Create.DryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")
	releaseSet.Create.Fs.Usage = func() {
		fmt.Println("Available options:")
		releaseSet.Create.Fs.PrintDefaults()
	}
	releaseSet.Publish.Fs = flag.NewFlagSet("release publish", flag.ExitOnError)
	releaseSet.Publish.Fs.IntVar(&releaseSet.Publish.PrNumber, "prnumber", 0, "overwrite the issue number")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.MergeCommitSHA, "merge-sha", "", "set the merge sha")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.ReleaseBranch, "release-branch", "", "set release branch (default: git default)")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.Version, "version", "", "override version to Update")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.Body, "body", "", "custom release message (markdow string or file)")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.PatchLevel, "patchLevel", "", "predefine version to Update")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.Assets, "assets", "", "upload assets to release. eg.: \"file=awesome-ci\"")
	releaseSet.Publish.Fs.Int64Var(&releaseSet.Publish.ReleaseId, "releaseid", 0, "publish an early defined release")
	releaseSet.Publish.Fs.BoolVar(&releaseSet.Publish.DryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")
	releaseSet.Publish.Fs.BoolVar(&releaseSet.Publish.Hotfix, "hotfix", false, "create a hotfix release")

	// parseJSON
	parseSet.Fs = flag.NewFlagSet("parse", flag.ExitOnError)
	parseSet.Fs.Usage = func() {
		fmt.Println("Available commands:")
		fmt.Println("  json")
		fmt.Println("  yaml")
		fmt.Println("Use \"awesome-ci release <command> --help\" for more information about a given command.")
	}
	parseSet.Json.Fs = flag.NewFlagSet("parse json", flag.ExitOnError)
	parseSet.Json.Fs.StringVar(&parseSet.Json.file, "file", "", "file to be parsed")
	parseSet.Json.Fs.StringVar(&parseSet.Json.value, "value", "", "value for output")
	parseSet.Json.Fs.Usage = func() {
		fmt.Println("Available options:")
		parseSet.Json.Fs.PrintDefaults()
	}
	parseSet.Yaml.Fs = flag.NewFlagSet("parse yaml", flag.ExitOnError)
	parseSet.Yaml.Fs.StringVar(&parseSet.Yaml.file, "file", "", "file to be parsed")
	parseSet.Yaml.Fs.StringVar(&parseSet.Yaml.value, "value", "", "value for output")
	parseSet.Yaml.Fs.Usage = func() {
		fmt.Println("Available options:")
		parseSet.Yaml.Fs.PrintDefaults()
	}
}

func returnHelpIfEmpty(args []string, usage func()) {
	if len(args) < 1 {
		usage()
		os.Exit(0)
	}
}

func printNoValidCommand(usage func()) {
	fmt.Println("The given command is not known!")
	fmt.Println("")
	usage()
	os.Exit(1)
}

func main() {

	flag.Usage = func() {
		fmt.Println("awesome-ci makes your CI easy.")
		fmt.Println("  Find more information and examples at: https://github.com/fullstack-devops/awesome-ci")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("  release")
		fmt.Println("  parse")
		fmt.Println("  pr")
		fmt.Println("")
		fmt.Println("Use \"awesome-ci <command> --help\" for more information about a given command.")
	}
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	returnHelpIfEmpty(flag.Args(), flag.Usage)

	// distribute environment settings
	/* environment := service.EvaluateEnvironment()
	service.CiEnvironment = environment */

	switch flag.Arg(0) {
	case "pr":
		pullrequestSet.Fs.Parse(flag.Args()[1:])
		returnHelpIfEmpty(flag.Args()[1:], pullrequestSet.Fs.Usage)
		switch flag.Arg(1) {
		case "info":
			err := pullrequestSet.Info.Fs.Parse(flag.Args()[2:])
			if err != nil {
				log.Fatalln(err)
			}
			service.PrintPRInfos(&pullrequestSet.Info)
		default:
			printNoValidCommand(pullrequestSet.Fs.Usage)
		}
	case "release":
		releaseSet.Fs.Parse(flag.Args()[1:])
		returnHelpIfEmpty(flag.Args()[1:], releaseSet.Fs.Usage)
		switch flag.Arg(1) {
		case "create":
			err := releaseSet.Create.Fs.Parse(flag.Args()[2:])
			if err != nil {
				log.Fatalln(err)
			}
			service.ReleaseCreate(&releaseSet.Create)
		case "publish":
			err := releaseSet.Publish.Fs.Parse(flag.Args()[2:])
			if err != nil {
				log.Fatalln(err)
			}
			service.ReleasePublish(&releaseSet.Publish)
		default:
			printNoValidCommand(releaseSet.Fs.Usage)
		}
	case "parse":
		parseSet.Fs.Parse(flag.Args()[1:])
		returnHelpIfEmpty(flag.Args()[1:], parseSet.Fs.Usage)
		switch flag.Arg(1) {
		case "json":
			parseSet.Json.Fs.Parse(flag.Args()[2:])
			service.ParseJson(&parseSet.Json.file, &parseSet.Json.value)
		case "yaml":
			parseSet.Yaml.Fs.Parse(flag.Args()[2:])
			service.ParseJson(&parseSet.Yaml.file, &parseSet.Yaml.value)
		default:
			printNoValidCommand(parseSet.Fs.Usage)
		}
	default:
		printNoValidCommand(flag.Usage)
	}
}
