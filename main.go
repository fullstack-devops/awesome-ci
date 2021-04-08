package main

import (
	"awesome-ci/service"
	"flag"
	"fmt"
<<<<<<< HEAD
=======
	"io/ioutil"
	"log"
>>>>>>> ec5109b3b8a4649530e5d2d430c0ee44624b9be8
	"os"
)

var (
	cienv         string
	createRelease CreateReleaseSet
	getBuildInfos getBuildInfosSet
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
	dryRun          bool
}

type getBuildInfosSet struct {
	fs         *flag.FlagSet
	version    string
	patchLevel string
	output     string
}

func init() {
<<<<<<< HEAD
	flag.StringVar(&cienv, "cienv", "Github", "set your CI Environment for Special Featueres!\nAvalible: Jenkins, Github, Gitlab, Custom\nDefault: Github")
=======
	cienv = flag.String("cienv", "Github", "set your CI Environment for Special Featueres!\nAvalible: Jenkins, Github, Gitlab, Custom\nDefault: Github")
>>>>>>> ec5109b3b8a4649530e5d2d430c0ee44624b9be8
	flag.BoolVar(&versionFlag, "version", false, "print version by calling it")
	// flag.BoolVar(&debug, "debug", false, "enable debug level by calling it")

	// createReleaseSet
	createRelease.fs = flag.NewFlagSet("createRelease", flag.ExitOnError)
	createRelease.fs.StringVar(&createRelease.version, "version", "", "override version to Update")
	createRelease.fs.StringVar(&createRelease.patchLevel, "patchLevel", "bugfix", "predefine version to Update")
	createRelease.fs.StringVar(&createRelease.publishNpm, "publishNpm", "", "runs npm publish --tag <createdTag> with custom directory")
	createRelease.fs.StringVar(&createRelease.uploadArtifacts, "uploadArtifacts", "", "uploads atifacts to release (file)")
	createRelease.fs.BoolVar(&createRelease.dryRun, "dry-run", false, "make dry-run before writing version to Git by calling it")

	// getNewReleaseVersion
	getBuildInfos.fs = flag.NewFlagSet("getBuildInfos", flag.ExitOnError)
	getBuildInfos.fs.StringVar(&getBuildInfos.version, "version", "", "override version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.patchLevel, "patchLevel", "bugfix", "predefine version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.output, "output", "pr,version", "define output by get")
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
		fmt.Print("\nUsage:\n  awesome-ci [subcommand] [options]\n")
		fmt.Print("\nUse awesome-ci createRelease -patchLevel bugfix -dry-run\n")
		fmt.Print("CI examples at: https://github.com/eksrvb/awesome-ci\n")
	}
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
	}

	switch os.Args[1] {
	case "createRelease":
		createRelease.fs.Parse(os.Args[2:])
		service.CreateRelease(cienv, &createRelease.version, &createRelease.patchLevel, &createRelease.dryRun, &createRelease.publishNpm, &createRelease.uploadArtifacts)
	case "getBuildInfos":
		getBuildInfos.fs.Parse(os.Args[2:])
		service.GetBuildInfos(cienv, &getBuildInfos.version, &getBuildInfos.patchLevel, &getBuildInfos.output)
	}
}
