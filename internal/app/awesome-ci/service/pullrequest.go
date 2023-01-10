package service

import (
	"awesome-ci/internal/app/awesome-ci/ces"
	scmportal "awesome-ci/internal/app/awesome-ci/scm-portal"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func PrintPRInfos(number int, formatOut string) {
	scmLayer, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		log.Fatalln(err)
	}

	if err = evalPrNumber(&number); err != nil {
		log.Fatalln(err)
	}

	prInfos, err := scmLayer.GetPrInfos(number, "")
	if err != nil {
		log.Fatalln(err)
	}

	var envVars []ces.KeyValue = []ces.KeyValue{
		{Name: "ACI_PR", Value: strconv.Itoa(prInfos.Number)},
		{Name: "ACI_PR_SHA", Value: prInfos.Sha},
		{Name: "ACI_PR_SHA_SHORT", Value: prInfos.ShaShort},
		{Name: "ACI_PR_BRANCH", Value: prInfos.BranchName},
		{Name: "ACI_MERGE_COMMIT_SHA", Value: prInfos.MergeCommitSha},
		{Name: "ACI_OWNER", Value: prInfos.Owner},
		{Name: "ACI_REPO", Value: prInfos.Repo},
		{Name: "ACI_PATCH_LEVEL", Value: string(prInfos.PatchLevel)},
		{Name: "ACI_VERSION", Value: prInfos.NextVersion},
		{Name: "ACI_LATEST_VERSION", Value: prInfos.LatestVersion},
	}

	if err := scmLayer.CES.ExportAsEnv(envVars); err != nil {
		log.Fatalln("could not export env variables: %v", err)
	}

	fmt.Println("### Info output:")
	fmt.Printf("Pull Request: %d\n", prInfos.Number)
	fmt.Printf("Latest release version: %s\n", prInfos.LatestVersion)
	fmt.Printf("Patch level: %s\n", prInfos.PatchLevel)
	fmt.Printf("Possible new release version: %s\n", prInfos.NextVersion)

}
