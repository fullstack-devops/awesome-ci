package service

import (
	"errors"
	"flag"
	"fmt"
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
	PublishNpm      string
	UploadArtifacts string
	PrNumber        int
	DryRun          bool
}

func ReleaseCreate(args *ReleaseCreateSet) {
	prNumber, err := evalPrNumber(&args.PrNumber)
	if err != nil {
		panic(err)
	}

	prInfos, err := getPRInfos(prNumber, true)
	if err != nil {
		panic(err)
	}

	if args.DryRun {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Would writing new release: %s\n", prInfos.NextVersion)
	} else {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Writing new release: %s\n", prInfos.NextVersion)
		CiEnvironment.CreateRelease(&prInfos, true)
	}
}

func ReleasePublish(args *ReleasePublishSet) {

	if args.PatchLevel != "" && args.Version != "" {

	}

	// TODO
	if args.DryRun {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Would publishing release: %s\n", prInfos.NextVersion)
	} else {
		fmt.Printf("Old version: %s\n", prInfos.CurrentVersion)
		fmt.Printf("Publishing release: %s\n", prInfos.NextVersion)
		CiEnvironment.PublishRelease(&prInfos, &args.UploadArtifacts)
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
