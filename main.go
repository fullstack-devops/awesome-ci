package main

import (
	"awesome-ci/src/ciRunnerController"
	"awesome-ci/src/gitOnlineController"
	"awesome-ci/src/service"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	cienv         string
	createRelease CreateReleaseSet
	getBuildInfos GetBuildInfosSet
	parseJson     parseJsonYamlSet
	parseYaml     parseJsonYamlSet
	version       string
	versionFlag   bool
	//debug         bool
)

type CreateReleaseSet struct {
	fs              *flag.FlagSet
	version         string
	patchLevel      string
	publishNpm      string
	uploadArtifacts string
	preRelease      bool
	dryRun          bool
}

type GetBuildInfosSet struct {
	fs         *flag.FlagSet
	version    string
	patchLevel string
	format     string
}

type parseJsonYamlSet struct {
	fs    *flag.FlagSet
	file  string
	value string
}

func init() {
	flag.StringVar(&cienv, "cienv", "", "set your CI Environment for Special Featueres!\nAvalible: Jenkins, Github, Gitlab, Custom\nDefault: Github")
	flag.BoolVar(&versionFlag, "version", false, "print version by calling it")
	// flag.BoolVar(&debug, "debug", false, "enable debug level by calling it")

	// createReleaseSet
	createRelease.fs = flag.NewFlagSet("createRelease", flag.ExitOnError)
	createRelease.fs.StringVar(&createRelease.version, "version", "", "override version to Update")
	createRelease.fs.StringVar(&createRelease.patchLevel, "patchLevel", "", "predefine version to Update")
	createRelease.fs.StringVar(&createRelease.publishNpm, "publishNpm", "", "runs npm publish --tag <createdTag> with custom directory")
	createRelease.fs.StringVar(&createRelease.uploadArtifacts, "uploadArtifacts", "", "uploads atifacts to release (file)")
	createRelease.fs.BoolVar(&createRelease.dryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")
	createRelease.fs.BoolVar(&createRelease.preRelease, "preRelease", false, "creates an github pre release")

	// getNewReleaseVersion
	getBuildInfos.fs = flag.NewFlagSet("getBuildInfos", flag.ExitOnError)
	getBuildInfos.fs.StringVar(&getBuildInfos.version, "version", "", "override version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.patchLevel, "patchLevel", "", "predefine version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.format, "format", "", "define output by get")

	// parseJSON
	parseJson.fs = flag.NewFlagSet("parseJSON", flag.ExitOnError)
	parseJson.fs.StringVar(&parseJson.file, "file", "", "file to be parsed")
	parseJson.fs.StringVar(&parseJson.value, "value", "", "value for output")

	// parseYAML
	parseYaml.fs = flag.NewFlagSet("parseYAML", flag.ExitOnError)
	parseYaml.fs.StringVar(&parseYaml.file, "file", "", "file to be parsed")
	parseYaml.fs.StringVar(&parseYaml.value, "value", "", "value for output")
}

func main() {
	flag.Usage = func() {
		fmt.Println("awesome-ci makes your CI easy.")
		fmt.Print("\n  Find more information and examples at: https://github.com/eksrvb/awesome-ci\n\n")
		fmt.Println("Available commands:")
		flag.PrintDefaults()
		fmt.Print("\nSubcommand: createRelease\n")
		createRelease.fs.PrintDefaults()
		fmt.Print("\nSubcommand: getBuildInfos\n")
		getBuildInfos.fs.PrintDefaults()
		fmt.Print("\nSubcommand: parseJSON, parseYAML\n")
		parseJson.fs.PrintDefaults()
		fmt.Print("\nUsage:\n  awesome-ci [subcommand] [options]\n")
		fmt.Print("\nUse awesome-ci createRelease -patchLevel bugfix -dry-run\n")
		fmt.Print("CI examples at: https://github.com/eksrvb/awesome-ci\n")
	}
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	// distribute environment settings
	environment := service.EvaluateEnvironment()
	ciRunnerController.CiEnvironment = environment
	gitOnlineController.CiEnvironment = environment

	switch flag.Args()[0] {
	case "createRelease":
		createRelease.fs.Parse(flag.Args()[1:])
		service.CreateRelease(cienv, &createRelease.version, &createRelease.patchLevel, &createRelease.dryRun, &createRelease.preRelease, &createRelease.publishNpm, &createRelease.uploadArtifacts)
	case "getBuildInfos":
		getBuildInfos.fs.Parse(flag.Args()[1:])
		service.GetBuildInfos(cienv, &getBuildInfos.version, &getBuildInfos.patchLevel, &getBuildInfos.format)
	case "parseJSON":
		parseJson.fs.Parse(flag.Args()[1:])
		service.ParseJson(&parseJson.file, &parseJson.value)
	case "parseYAML":
		parseYaml.fs.Parse(flag.Args()[1:])
		service.ParseYaml(&parseYaml.file, &parseYaml.value)
	default:
		log.Fatalln("Not a valid Subcommand")
	}
}
