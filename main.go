package main

import (
	"awesome-ci/service"
	"flag"
	"os"
)

var (
	cienv         *string
	createRelease CreateReleaseSet
	getBuildInfos getBuildInfosSet
)

type CreateReleaseSet struct {
	fs         *flag.FlagSet
	version    string
	patchLevel string
	publishNpm string
	dryRun     bool
}

type getBuildInfosSet struct {
	fs         *flag.FlagSet
	version    string
	patchLevel string
	output     string
}

func init() {
	cienv = flag.String("cienv", "Github", "set your CI Environment for Special Featueres!\nAvalible: Jenkins, Github, Gitlab, Custom\nDefault: Github")

	// createReleaseSet
	createRelease.fs = flag.NewFlagSet("createRelease", flag.ExitOnError)
	createRelease.fs.StringVar(&createRelease.version, "version", "", "override version to Update")
	createRelease.fs.StringVar(&createRelease.patchLevel, "patchLevel", "bugfix", "predefine version to Update")
	createRelease.fs.StringVar(&createRelease.publishNpm, "publishNpm", "", "runs npm publish --tag <createdTag> with custom directory")
	createRelease.fs.BoolVar(&createRelease.dryRun, "dry-run", false, "make dry-run before writing version to Git")

	// getNewReleaseVersion
	getBuildInfos.fs = flag.NewFlagSet("getBuildInfos", flag.ExitOnError)
	getBuildInfos.fs.StringVar(&getBuildInfos.version, "version", "", "override version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.patchLevel, "patchLevel", "bugfix", "predefine version to Update")
	getBuildInfos.fs.StringVar(&getBuildInfos.output, "output", "pr,version", "define output by get")
}

func main() {

	// disable logging
	//log.SetOutput(ioutil.Discard)

	flag.Parse()

	switch os.Args[1] {
	case "createRelease":
		createRelease.fs.Parse(os.Args[2:])
		service.CreateRelease(cienv, &createRelease.version, &createRelease.patchLevel, &createRelease.dryRun, &createRelease.publishNpm)
	case "getBuildInfos":
		getBuildInfos.fs.Parse(os.Args[2:])
		service.GetBuildInfos(cienv, &getBuildInfos.version, &getBuildInfos.patchLevel, &getBuildInfos.output)
	}

}
