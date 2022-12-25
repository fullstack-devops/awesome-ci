package service

import (
	"awesome-ci/internal/app/awesome-ci/acigithub"
	"awesome-ci/internal/app/awesome-ci/models"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runcmd(cmd string, shell bool) string {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println(err)
		}
		return string(out)
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

func standardPrInfosToEnv(prInfos *models.StandardPrInfos) (err error) {
	runnerType := "github_runner"
	switch runnerType {
	case "github_runner":
		envVars, err := acigithub.OpenEnvFile()
		if err != nil {
			return err
		}
		envVars.Set("ACI_PR", strconv.Itoa(prInfos.PrNumber))
		envVars.Set("ACI_PR_SHA", prInfos.Sha)
		envVars.Set("ACI_PR_SHA_SHORT", prInfos.ShaShort)
		envVars.Set("ACI_PR_BRANCH", prInfos.BranchName)
		envVars.Set("ACI_MERGE_COMMIT_SHA", prInfos.MergeCommitSha)

		envVars.Set("ACI_OWNER", prInfos.Owner)
		envVars.Set("ACI_REPO", prInfos.Repo)
		envVars.Set("ACI_PATCH_LEVEL", string(prInfos.PatchLevel))
		envVars.Set("ACI_VERSION", prInfos.NextVersion)
		envVars.Set("ACI_LATEST_VERSION", prInfos.LatestVersion)

		err = envVars.SaveEnvFile()
		if err != nil {
			return err
		}
	default:
		log.Println("Runner Type not implemented!")
	}
	return
}
