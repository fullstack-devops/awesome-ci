package service

import (
	scmportal "awesome-ci/internal/app/awesome-ci/scm-portal"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func PrintPRInfos(number int, format string) {
	cesType, grc, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		log.Fatalln(err)
	}

	if err = evalPrNumber(&number); err != nil {
		log.Fatalln(err)
	}

	prInfos, err := scmportal.GetPrInfos(cesType, grc, number, "")
	if err != nil {
		log.Fatalln(err)
	}

	/* 	if errEnvs := standardPrInfosToEnv(prInfos); errEnvs != nil {
		log.Fatalln(errEnvs)
	} */

	if format != "" {
		replacer := strings.NewReplacer(
			"pr", fmt.Sprint(prInfos.Number),
			"version", prInfos.NextVersion,
			"latest_version", prInfos.LatestVersion,
			"patchLevel", string(prInfos.PatchLevel))
		output := replacer.Replace(format)
		fmt.Print(output)
	} else {
		fmt.Println("#### Info output:")
		fmt.Printf("Pull Request: %d\n", prInfos.Number)
		fmt.Printf("Latest release version: %s\n", prInfos.LatestVersion)
		fmt.Printf("Patch level: %s\n", prInfos.PatchLevel)
		fmt.Printf("Possible new release version: %s\n", prInfos.NextVersion)
	}
}
