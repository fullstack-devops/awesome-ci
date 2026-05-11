package service

import (
	"fmt"
	"strconv"

	"github.com/layer87-labs/relctl/internal/app/relctl/ces"
	scmportal "github.com/layer87-labs/relctl/internal/app/relctl/scm-portal"

	log "github.com/sirupsen/logrus"
)

func PrintPRInfos(number int, mergeCommitSha string, formatOut string) {
	scmLayer, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof("detected ces type: %s", scmLayer.CES.Type)

	if mergeCommitSha == "" {
		if err = evalPrNumber(&number); err != nil {
			log.Fatalln(err)
		}
		log.Infof("evaluated pull request number %d", number)
	}

	prInfos, err := scmLayer.GetPrInfos(number, mergeCommitSha)
	if err != nil {
		log.Fatalln(err)
	}

	if err := prInfosToEnv(scmLayer, prInfos); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("### Info output:")
	fmt.Printf("Pull Request: %d\n", prInfos.Number)
	fmt.Printf("Latest release version: %s\n", prInfos.LatestVersion)
	fmt.Printf("Patch level: %s\n", prInfos.PatchLevel)
	fmt.Printf("Possible new release version: %s\n", prInfos.NextVersion)
}

func prInfosToEnv(scmLayer *scmportal.SCMLayer, prInfos *scmportal.PrMrRequestInfos) error {
	var envVars = []ces.KeyValue{
		{Name: "RELCTL_PR", Value: strconv.Itoa(prInfos.Number)},
		{Name: "RELCTL_PR_SHA", Value: prInfos.Sha},
		{Name: "RELCTL_PR_SHA_SHORT", Value: prInfos.ShaShort},
		{Name: "RELCTL_PR_BRANCH", Value: prInfos.BranchName},
		{Name: "RELCTL_MERGE_COMMIT_SHA", Value: prInfos.MergeCommitSha},
		{Name: "RELCTL_OWNER", Value: prInfos.Owner},
		{Name: "RELCTL_REPO", Value: prInfos.Repo},
		{Name: "RELCTL_PATCH_LEVEL", Value: string(prInfos.PatchLevel)},
		{Name: "RELCTL_VERSION", Value: prInfos.NextVersion},
		{Name: "RELCTL_NEXT_VERSION", Value: prInfos.NextVersion},
		{Name: "RELCTL_LATEST_VERSION", Value: prInfos.LatestVersion},
	}

	if err := scmLayer.CES.ExportAsEnv(envVars); err != nil {
		return fmt.Errorf("could not export env variables: %v", err)
	}
	return nil
}
