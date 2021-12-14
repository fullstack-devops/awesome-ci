package service

import (
	"awesome-ci/src/controlEnvs"
	"awesome-ci/src/gitController"
	"flag"
	"fmt"
	"log"
	"strconv"
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
	prInfos, err := getPRInfos(args.Number, true)
	if err != nil {
		log.Fatalln(err)
	}
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
	}
}

func getPRInfos(prNumber int, silent bool) (aciPrInfos gitController.AciPrInfos, err error) {
	aciPrInfos, err = CiEnvironment.GetPrInfos(prNumber)
	switch CiEnvironment.RunnerType {
	case "github_runner":
		envVariables, err := controlEnvs.OpenEnvFile(CiEnvironment.RunnerInfo.EnvFile)
		if err != nil {
			return aciPrInfos, err
		}
		envVariables.Set("ACI_PR", strconv.Itoa(prNumber))
		envVariables.Set("ACI_PR_SHA", aciPrInfos.Sha)
		envVariables.Set("ACI_PR_SHA_SHORT", aciPrInfos.ShaShort)
		envVariables.Set("ACI_PR_BRANCH", strings.ToLower(*CiEnvironment.GitInfos.Owner))
		envVariables.Set("ACI_ORGA", strings.ToLower(*CiEnvironment.GitInfos.Repo))
		envVariables.Set("ACI_PATCH_LEVEL", aciPrInfos.PatchLevel)
		envVariables.Set("ACI_VERSION", aciPrInfos.NextVersion)
		envVariables.Set("ACI_LATEST_VERSION", aciPrInfos.LatestVersion)
		err = envVariables.SaveEnvFile()
		if err != nil {
			return aciPrInfos, err
		}
	}
	return
}
