package service

import (
	"awesome-ci/internal/app/awesome-ci/acigithub"
	"awesome-ci/internal/app/awesome-ci/semver"
	"awesome-ci/internal/app/awesome-ci/tools"
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/google/go-github/v44/github"
)

type ReleaseSet struct {
	Fs      *flag.FlagSet
	Create  ReleaseCreateSet
	Publish ReleasePublishSet
}

type ReleaseCreateSet struct {
	Fs             *flag.FlagSet
	Version        string
	PatchLevel     string
	PrNumber       int
	MergeCommitSHA string
	ReleaseBranch  string
	DryRun         bool
	Hotfix         bool
	Body           string
}

type ReleasePublishSet struct {
	Fs             *flag.FlagSet
	Version        string
	PatchLevel     string
	ReleaseId      int64
	Assets         string
	PrNumber       int
	MergeCommitSHA string
	ReleaseBranch  string
	DryRun         bool
	Hotfix         bool
	Body           string
}

func ReleaseCreate(args *ReleaseCreateSet) *github.RepositoryRelease {
	var version string = ""
	_, err := acigithub.NewGitHubClient()
	if err != nil {
		log.Fatalln(err)
	}

	if args.Version != "" && args.PatchLevel != "" {
		parsedPatchLevel := semver.ParsePatchLevel(args.PatchLevel)
		version, err = semver.IncreaseVersion(parsedPatchLevel, args.Version)
		if err != nil {
			log.Fatalln(err)
		}
	} else if args.Version != "" && args.PatchLevel == "" {
		version = args.Version
	} else if args.Hotfix {

		release, err := acigithub.GetLatestReleaseVersion()

		if err != nil {
			log.Fatalln(err)
		}
		version, err = semver.IncreaseVersion(semver.Bugfix, *release.TagName)

		if err != nil {
			log.Fatalln(err)
		}

	} else {
		// if no merge commit sha is provided, the pull request number should either be specified or evaluated from the merge message (fallback)
		if args.MergeCommitSHA == "" {
			err := evalPrNumber(&args.PrNumber)
			if err != nil {
				log.Fatalln(err)
			}
		}

		prInfos, _, err := acigithub.GetPrInfos(args.PrNumber, args.MergeCommitSHA)
		if err != nil {
			log.Fatalln(err)
		}

		version = prInfos.NextVersion
		if errEnvs := standardPrInfosToEnv(prInfos); errEnvs != nil {
			log.Fatalln(errEnvs)
		}
	}

	if args.DryRun {
		fmt.Printf("Would create new release with version: %s\n", version)
	} else {
		fmt.Printf("Writing new release: %s\n", version)
		createdRelease, err := acigithub.CreateRelease(version, args.ReleaseBranch, args.Body, true)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Create release successful. ID:", *createdRelease.ID)

		return createdRelease
	}

	return nil
}

func ReleasePublish(args *ReleasePublishSet) {
	var version string = ""
	_, err := acigithub.NewGitHubClient()
	if err != nil {
		log.Fatalln(err)
	}

	if args.Version != "" && args.PatchLevel != "" {
		parsedPatchLevel := semver.ParsePatchLevel(args.PatchLevel)
		version, err = semver.IncreaseVersion(parsedPatchLevel, args.Version)
		if err != nil {
			log.Fatalln(err)
		}
	} else if args.Version != "" && args.PatchLevel == "" {
		version = args.Version
	} else if args.Hotfix {

		release, err := acigithub.GetLatestReleaseVersion()

		if err != nil {
			log.Fatalln(err)
		}
		version, err = semver.IncreaseVersion(semver.Bugfix, *release.TagName)

		if err != nil {
			log.Fatalln(err)
		}

	} else if args.ReleaseId == 0 {
		// if no merge commit sha is provided, the pull request number should either be specified or evaluated from the merge message (fallback)
		if args.MergeCommitSHA == "" {
			err := evalPrNumber(&args.PrNumber)
			if err != nil {
				log.Fatalln(err)
			}
		}
		prInfos, _, err := acigithub.GetPrInfos(args.PrNumber, args.MergeCommitSHA)
		if err != nil {
			log.Fatalln(err)
		}
		version = prInfos.NextVersion
		if errEnvs := standardPrInfosToEnv(prInfos); errEnvs != nil {
			log.Fatalln(errEnvs)
		}
	}

	if args.Assets != "" {
		_, err = tools.GetAsstes(&args.Assets, true)
		if err != nil {
			log.Fatalln("not all specified asstes available, please check")
		}
	}

	if args.DryRun {
		fmt.Printf("Would publishing release: %s\n", version)
	} else {
		fmt.Printf("Publishing release: %s - %d\n", version, args.ReleaseId)
		err = acigithub.PublishRelease(version, args.ReleaseBranch, args.Body, args.ReleaseId, &args.Assets)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func evalPrNumber(override *int) (err error) {
	if *override != 0 {
		return nil
	}

	*override, err = getPrFromMergeMessage()
	if err != nil {
		return err
	}
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
			"\nAlternativly provide the PR-Number by adding the argument -number <int>")
	}
}
