package service

import (
	"awesome-ci/internal/pkg/detect"
	"awesome-ci/internal/pkg/models"
	"fmt"
	"os/exec"
	"strconv"
)

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
	envs, err := detect.LoadEnvVars()
	if err != nil {
		return err
	}
	envs.Set("ACI_PR", strconv.Itoa(prInfos.PrNumber))
	envs.Set("ACI_PR_SHA", prInfos.Sha)
	envs.Set("ACI_PR_SHA_SHORT", prInfos.ShaShort)
	envs.Set("ACI_PR_BRANCH", prInfos.BranchName)
	envs.Set("ACI_MERGE_COMMIT_SHA", prInfos.MergeCommitSha)

	envs.Set("ACI_OWNER", prInfos.Owner)
	envs.Set("ACI_REPO", prInfos.Repo)
	envs.Set("ACI_PATCH_LEVEL", string(prInfos.PatchLevel))
	envs.Set("ACI_VERSION", prInfos.NextVersion)
	envs.Set("ACI_LATEST_VERSION", prInfos.LatestVersion)

	envs.SetEnvVars()

	return nil
}
