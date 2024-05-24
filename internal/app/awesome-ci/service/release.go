package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces"
	scmportal "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	log "github.com/sirupsen/logrus"
)

type ReleaseArgs struct {
	Version        string
	PatchLevel     string
	PrNumber       int
	MergeCommitSHA string
	ReleaseBranch  string
	ReleasePrefix  string
	DryRun         bool
	Hotfix         bool
	Body           string
}

func ReleaseCreate(args *ReleaseArgs) {
	scmLayer, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		log.Fatalln(err)
	}

	version, releasePrefix, err := argsToVersion(scmLayer, args)
	if err != nil {
		log.Fatalln(err)
	}

	if args.DryRun {
		log.Infof("Would create new release with version: %s\n", version)

		fmt.Println("### Info output:")
		fmt.Printf("Would create new release with version: %s\n", version)
	} else {
		log.Infof("Writing new release: %s\n", version)
		fmt.Println("### Info output:")
		fmt.Printf("Writing new release with version: %s\n", version)

		createdRelease, err := scmLayer.CreateRelease(version, releasePrefix, args.ReleaseBranch, args.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("Create release successful. ID: %d", createdRelease.ID)
		fmt.Printf("Create release successful. ID: %d", createdRelease.ID)

		var envVars = []ces.KeyValue{
			{Name: "ACI_RELEASE_ID", Value: fmt.Sprintf("%d", createdRelease.ID)},
		}

		if err := scmLayer.CES.ExportAsEnv(envVars); err != nil {
			log.Fatalf("could not export env variables: %v", err)
		}
	}
}

func ReleasePublish(args *ReleaseArgs, releaseID int64, assets []string) {
	scmLayer, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		log.Fatalln(err)
	}

	var version, releasePrefix string

	if releaseID == 0 {
		version, releasePrefix, err = argsToVersion(scmLayer, args)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var assetsEncoded []tools.UploadAsset
	if len(assets) > 0 {
		for _, asset := range assets {
			assetInfo, err := tools.GetAsset(asset)
			if err != nil {
				log.Fatalln("not all specified assets available, please check", err)
			}
			assetsEncoded = append(assetsEncoded, *assetInfo)
		}
	}
	log.Infof("will upload %d assets", len(assets))

	body, err := tools.ReadFileToString(args.Body)
	if err != nil {
		log.Warnf("could not process the given body: %v", err)
	}

	if args.DryRun {
		log.Infof("Would publishing release: %s", version)
	} else {
		log.Infof("Publishing release: %s - %d", version, releaseID)
		_, err := scmLayer.PublishRelease(version, releasePrefix, args.ReleaseBranch, body, releaseID, assetsEncoded)
		if err != nil {
			log.Fatalln(err)
		}

		/* var envVars []ces.KeyValue
		for i, ra := range relAssets {
			// export Download URL to env. See: #53
			envVars = append(envVars, ces.KeyValue{Name: fmt.Sprintf("ACI_ARTIFACT_%d_URL", i+1), Value: *ra.BrowserDownloadURL})
		}

		if len(envVars) > 0 {
			if err := scmLayer.CES.ExportAsEnv(envVars); err != nil {
				log.Fatalln("could not export env variables: %v", err)
			}
		} */
	}
}

func argsToVersion(scmLayer *scmportal.SCMLayer, args *ReleaseArgs) (version, releasePrefix string, err error) {
	if args.ReleasePrefix != "" {
		releasePrefix = args.ReleasePrefix
	} else {
		releasePrefix = "Release"
	}

	if args.Version != "" && args.PatchLevel != "" {
		parsedPatchLevel, err := semver.ParsePatchLevel(args.PatchLevel)
		if err != nil && err != semver.ErrUseMinimalPatchVersion {
			log.Fatalln(err)
		}
		version, err = semver.IncreaseVersion(parsedPatchLevel, args.Version)
		if err != nil {
			log.Fatalln(err)
		}

	} else if args.Version != "" && args.PatchLevel == "" {
		version = args.Version

	} else if args.Hotfix {
		release, err := scmLayer.GetLatestReleaseVersion()
		if err != nil {
			log.Fatalln(err)
		}

		version, err = semver.IncreaseVersion(semver.Bugfix, release.TagName)
		if err != nil {
			log.Fatalln(err)
		}

		if args.ReleasePrefix == "" {
			releasePrefix = "Hotfix"
		}

	} else {
		// if no merge commit sha is provided, the pull request number should either be specified or evaluated from the merge message (fallback)
		if args.MergeCommitSHA == "" {
			log.Infoln("no merge commit sha is given, eval...")
			err := evalPrNumber(&args.PrNumber)
			if err != nil {
				log.Fatalln(err)
			}
		}

		prInfos, err := scmLayer.GetPrInfos(args.PrNumber, args.MergeCommitSHA)
		if err != nil {
			log.Fatalln(err)
		}

		version = prInfos.NextVersion

		if err = prInfosToEnv(scmLayer, prInfos); err != nil {
			log.Fatalln(err)
		}
	}
	return
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

	mergeMessage := r.FindStringSubmatch(tools.RunCmd(`git log -1 --pretty=format:"%s"`, true))
	if len(mergeMessage) > 1 {
		return strconv.Atoi(mergeMessage[1])
	} else {
		return 0, errors.New("No PR found in merge message pls make shure this regex matches: " + regex +
			"\nExample: Merge pull request #3 from some-orga/feature/awesome-feature" +
			"\nAlternativly provide the PR-Number by adding the argument -number <int>")
	}
}
