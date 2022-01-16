package service

import (
	"awesome-ci/src/acigithub"
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type ReleaseSet struct {
	Fs      *flag.FlagSet
	Create  ReleaseCreateSet
	Publish ReleasePublishSet
}

type ReleaseCreateSet struct {
	Fs         *flag.FlagSet
	Version    string
	PatchLevel string
	PrNumber   int
	DryRun     bool
}

type ReleasePublishSet struct {
	Fs              *flag.FlagSet
	Version         string
	PatchLevel      string
	ReleaseId       string
	PublishNpm      string
	UploadArtifacts string
	PrNumber        int
	DryRun          bool
}

func ReleaseCreate(args *ReleaseCreateSet) {
	_, err := acigithub.NewGitHubClient()
	if err != nil {
		log.Fatalln(err)
	}
	prNumber, err := evalPrNumber(&args.PrNumber)
	if err != nil {
		log.Fatalln(err)
	}

	prInfos, _, err := acigithub.GetPrInfos(prNumber)
	if err != nil {
		log.Fatalln(err)
	}

	if args.DryRun {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Would create new release: %s\n", prInfos.NextVersion)
	} else {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Writing new release: %s\n", prInfos.NextVersion)
		createdRelease, err := acigithub.CreateRelease(prInfos, true)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Create release successful. ID:", *createdRelease.ID)
		envVars, err := acigithub.OpenEnvFile()
		if err != nil {
			log.Fatalln(err)
		}
		envVars.Set("ACI_RELEASE_ID", fmt.Sprint(*createdRelease.ID))
		envVars.SaveEnvFile()
	}
}

func ReleasePublish(args *ReleasePublishSet) {
	_, err := acigithub.NewGitHubClient()
	if err != nil {
		log.Fatalln(err)
	}
	prNumber, err := evalPrNumber(&args.PrNumber)
	if err != nil {
		log.Fatalln(err)
	}

	prInfos, _, err := acigithub.GetPrInfos(prNumber)
	if err != nil {
		log.Fatalln(err)
	}
	/* if args.PatchLevel != "" && args.Version != "" {

	} */

	if args.DryRun {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Would publishing release: %s\n", prInfos.NextVersion)
	} else {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Publishing release: %s\n", prInfos.NextVersion)
		acigithub.PublishRelease(prInfos, &args.UploadArtifacts)
	}
}

func evalPrNumber(override *int) (prNumber int, err error) {
	if *override != 0 {
		return *override, nil
	}

	prNumber, err = getPrFromMergeMessage()
	if err != nil {
		panic(err)
	}

	/* if prNumber == 0 {
		// tags, _, _ := CiEnvironment.Clients.GithubClient.Repositories.ListTags()
	} */
	return
}

func getPrFromMergeMessage() (pr int, err error) {
	regex := `.*#([0-9]+).*`
	r := regexp.MustCompile(regex)

	mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
	if len(mergeMessage) > 1 {
		return strconv.Atoi(mergeMessage[1])
	} else {
		return 0, errors.New("No PR found in merge message pls make shure this regex matches: " + regex +
			"\nExample: Merge pull request #3 from some-orga/feature/awesome-feature" +
			"\nIf you like to set your patch level manually by flag: -level (feautre|bugfix)")
	}
}
