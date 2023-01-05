package service

import (
	"awesome-ci/internal/app/awesome-ci/connect"
	"awesome-ci/internal/pkg/detect"
	"awesome-ci/internal/pkg/semver"
	"awesome-ci/internal/pkg/tools"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v49/github"
)

type ReleaseCreateSet struct {
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
	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		log.Fatalln(err)
	}

	var version string = ""

	if args.Version != "" && args.PatchLevel != "" {
		parsedPatchLevel := semver.ParsePatchLevel(args.PatchLevel)
		version, err = semver.IncreaseVersion(parsedPatchLevel, args.Version)
		if err != nil {
			log.Fatalln(err)
		}
	} else if args.Version != "" && args.PatchLevel == "" {
		version = args.Version
	} else if args.Hotfix {

		release, err := ghrc.GetLatestReleaseVersion()

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

		prInfos, _, err := ghrc.GetPrInfos(args.PrNumber, args.MergeCommitSHA)
		if err != nil {
			log.Fatalln(err)
		}

		version = prInfos.NextVersion
		if errEnvs := standardPrInfosToEnv(prInfos); errEnvs != nil {
			log.Fatalln(errEnvs)
		}
	}

	if args.DryRun {
		log.Infof("Would create new release with version: %s\n", version)
	} else {
		log.Infof("Writing new release: %s\n", version)
		createdRelease, err := ghrc.CreateRelease(version, args.ReleaseBranch, args.Body, true)
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("Create release successful. ID: %s", *createdRelease.ID)

		envs, err := detect.LoadEnvVars()
		if err != nil {
			log.Warnf("could load env variables: %v", err)
		}
		envs.Set("ACI_RELEASE_ID", fmt.Sprintf("%d", *createdRelease.ID))
		if errEnvs := envs.SetEnvVars(); errEnvs != nil {
			log.Warnf("could not export env variable ACI_RELEASE_ID: %v", err)
		}

		return createdRelease
	}

	return nil
}

func ReleasePublish(args *ReleasePublishSet) {
	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		log.Fatalln(err)
	}

	var version string = ""

	if args.Version != "" && args.PatchLevel != "" {
		parsedPatchLevel := semver.ParsePatchLevel(args.PatchLevel)
		version, err = semver.IncreaseVersion(parsedPatchLevel, args.Version)
		if err != nil {
			log.Fatalln(err)
		}
	} else if args.Version != "" && args.PatchLevel == "" {
		version = args.Version
	} else if args.Hotfix {

		release, err := ghrc.GetLatestReleaseVersion()

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
		prInfos, _, err := ghrc.GetPrInfos(args.PrNumber, args.MergeCommitSHA)
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
		log.Infof("Would publishing release: %s\n", version)
	} else {
		log.Infof("Publishing release: %s - %d\n", version, args.ReleaseId)
		relAssets, err := ghrc.PublishRelease(version, args.ReleaseBranch, args.Body, args.ReleaseId, &args.Assets)
		if err != nil {
			log.Fatalln(err)
		}
		for i, ra := range relAssets {
			// export Download URL to env. See: #53
			envVars, err := detect.LoadEnvVars()
			if err != nil {
				log.Warnf("could load env variables: %v", err)
			}
			envVars.Set(fmt.Sprintf("ACI_ARTIFACT_%d_URL", i+1), *ra.BrowserDownloadURL)
			err = envVars.SetEnvVars()
			if err != nil {
				log.Warnf("could not export env variable ACI_RELEASE_ID: %v", err)
			}
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
