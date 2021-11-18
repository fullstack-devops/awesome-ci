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
	releaseSet     ReleaseSet
	pullrequestSet service.PullRequestSet
	parseSet       ParseSet
	version        string
	versionFlag    bool
)

type ReleaseSet struct {
	Fs     *flag.FlagSet
	Create struct {
		Fs         *flag.FlagSet
		version    string
		patchLevel string
		dryRun     bool
	}
	Publish struct {
		Fs              *flag.FlagSet
		version         string
		patchLevel      string
		publishNpm      string
		uploadArtifacts string
		dryRun          bool
	}
}

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
	pullrequestSet.Fs = flag.NewFlagSet("pullrequest", flag.ExitOnError)
	pullrequestSet.Fs.Usage = func() {
		fmt.Println("Available commands:")
		fmt.Println("  info")
		fmt.Println("Use \"awesome-ci pullrequest <command> --help\" for more information about a given command.")
	}
	pullrequestSet.Info.Fs = flag.NewFlagSet("pullrequest info", flag.ExitOnError)
	pullrequestSet.Info.Fs.IntVar(&pullrequestSet.Info.Number, "number", 0, "overwrite the issue number")
	pullrequestSet.Info.Fs.BoolVar(&pullrequestSet.Info.EvalNumber, "eval", false, "indicates whether the number should be determined automatically")

	// ReleaseSet
	releaseSet.Fs = flag.NewFlagSet("release", flag.ExitOnError)
	releaseSet.Fs.Usage = func() {
		fmt.Println("Available commands:")
		fmt.Println("  info")
		fmt.Println("  create")
		fmt.Println("  publish")
		fmt.Println("Use \"awesome-ci release <command> --help\" for more information about a given command.")
	}
	releaseSet.Create.Fs = flag.NewFlagSet("release create", flag.ExitOnError)
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.version, "version", "", "override version to Update")
	releaseSet.Create.Fs.StringVar(&releaseSet.Create.patchLevel, "patchLevel", "", "predefine version to Update")
	releaseSet.Create.Fs.BoolVar(&releaseSet.Create.dryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")
	releaseSet.Create.Fs.Usage = func() {
		fmt.Println("Available options:")
		releaseSet.Create.Fs.PrintDefaults()
	}
	releaseSet.Publish.Fs = flag.NewFlagSet("release publish", flag.ExitOnError)
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.version, "version", "", "override version to Update")
	releaseSet.Publish.Fs.StringVar(&releaseSet.Publish.patchLevel, "patchLevel", "", "predefine version to Update")
	releaseSet.Publish.Fs.BoolVar(&releaseSet.Publish.dryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")

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
		fmt.Println("  Find more information and examples at: https://github.com/eksrvb/awesome-ci")
		fmt.Println()
		fmt.Println("Available commands:")
		fmt.Println("  pullrequest")
		fmt.Println("  release")
		fmt.Println("  parse")
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
	environment := service.EvaluateEnvironment()
	service.CiEnvironment = environment

	switch flag.Arg(0) {
	case "pullrequest":
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
			service.ReleaseCreate(&releaseSet.Create.version, &releaseSet.Create.patchLevel, &releaseSet.Create.dryRun)
		case "publish":
			err := releaseSet.Publish.Fs.Parse(flag.Args()[2:])
			if err != nil {
				log.Fatalln(err)
			}
		default:
			printNoValidCommand(releaseSet.Fs.Usage)
		}
	case "getBuildInfos":
		// getBuildInfos.fs.Parse(flag.Args()[1:])
		// service.GetBuildInfos(cienv, &getBuildInfos.version, &getBuildInfos.patchLevel, &getBuildInfos.format)
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
