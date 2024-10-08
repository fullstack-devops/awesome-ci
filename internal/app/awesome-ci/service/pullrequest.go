package service

import (
	"fmt"
	"strconv"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/ces"
	scmportal "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal"

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
		{Name: "ACI_PR", Value: strconv.Itoa(prInfos.Number)},
		{Name: "ACI_PR_SHA", Value: prInfos.Sha},
		{Name: "ACI_PR_SHA_SHORT", Value: prInfos.ShaShort},
		{Name: "ACI_PR_BRANCH", Value: prInfos.BranchName},
		{Name: "ACI_MERGE_COMMIT_SHA", Value: prInfos.MergeCommitSha},
		{Name: "ACI_OWNER", Value: prInfos.Owner},
		{Name: "ACI_REPO", Value: prInfos.Repo},
		{Name: "ACI_PATCH_LEVEL", Value: string(prInfos.PatchLevel)},
		{Name: "ACI_VERSION", Value: prInfos.NextVersion},
		{Name: "ACI_NEXT_VERSION", Value: prInfos.NextVersion},
		{Name: "ACI_LATEST_VERSION", Value: prInfos.LatestVersion},
	}

	if err := scmLayer.CES.ExportAsEnv(envVars); err != nil {
		return fmt.Errorf("could not export env variables: %v", err)
	}
	return nil
}
