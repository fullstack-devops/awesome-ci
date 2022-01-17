package service

import (
	"awesome-ci/src/acigithub"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type PullRequestSet struct {
	Fs   *flag.FlagSet
	Info PullRequestInfoSet
}

type PullRequestInfoSet struct {
	Fs     *flag.FlagSet
	Number int
	Format string
}

func PrintPRInfos(args *PullRequestInfoSet) {
	_, err := acigithub.NewGitHubClient()
	if err != nil {
		log.Fatalln(err)
	}
	prInfos, _, err := acigithub.GetPrInfos(args.Number)
	if err != nil {
		log.Fatalln(err)
	}

	errEnvs := standardPrInfosToEnv(prInfos)
	if args.Format != "" {
		replacer := strings.NewReplacer(
			"pr", fmt.Sprint(prInfos.PrNumber),
			"version", prInfos.NextVersion,
			"latest_version", prInfos.LatestVersion,
			"patchLevel", prInfos.PatchLevel)
		output := replacer.Replace(args.Format)
		fmt.Print(output)
	} else {
		fmt.Println("#### Info output:")
		fmt.Printf("Pull Request: %d\n", prInfos.PrNumber)
		fmt.Printf("Latest release version: %s\n", prInfos.LatestVersion)
		fmt.Printf("Patch level: %s\n", prInfos.PatchLevel)
		fmt.Printf("Possible new release version: %s\n", prInfos.NextVersion)
		if errEnvs != nil {
			log.Fatalln(errEnvs)
		}
	}
	os.Exit(0)
}
