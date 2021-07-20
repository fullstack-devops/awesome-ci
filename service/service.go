package service

import (
	"awesome-ci/models"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var CiEnvironment models.CIEnvironment

type infosMergeMessage struct {
	PRNumber   string
	PatchLevel string
}

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

func EvaluateEnvironment() models.CIEnvironment {
	// env to check for github_runner
	githubRunnerApi, githubRunnerApiBool := os.LookupEnv("GITHUB_API_URL")
	githubRunnerRep, githubRunnerRepBool := os.LookupEnv("GITHUB_REPOSITORY")
	githubEnv, githubEnvBool := os.LookupEnv("GITHUB_ENV")
	// env to check for jenkins pipelines
	gitlabCi, gitlabCiBool := os.LookupEnv("GITLAB_CI")
	// env to check for jenkins pipelines
	_, jenkinsUrlBool := os.LookupEnv("JENKINS_URL")

	if githubRunnerApiBool && githubRunnerRepBool && githubEnvBool {
		CiEnvironment.GitType = "github"

		if !strings.HasSuffix(githubRunnerApi, "/") {
			githubRunnerApi = githubRunnerApi + "/"
		}
		CiEnvironment.GitInfos.ApiUrl = githubRunnerApi

		CiEnvironment.GitInfos.FullRepo = githubRunnerRep
		CiEnvironment.GitInfos.Orga = strings.Split(githubRunnerRep, "/")[0]
		CiEnvironment.GitInfos.Repo = strings.Split(githubRunnerRep, "/")[1]
		githubRunnerToken, githubRunnerTokenBool := os.LookupEnv("GITHUB_TOKEN")
		if !githubRunnerTokenBool {
			log.Fatalln("Apparently you are using a GitHub-Runner.\nPlease provide the GITHUB_TOKEN!\nSee https://docs.github.com/en/actions/reference/authentication-in-a-workflow#using-the-github_token-in-a-workflow\nand https://eksrvb.github.io/awesome-ci/examples/github_actions.html")
		}
		CiEnvironment.GitInfos.ApiToken = githubRunnerToken

		CiEnvironment.RunnerType = "github_runner"
		CiEnvironment.RunnerInfo.EnvFile = githubEnv

		return CiEnvironment
	} else if jenkinsUrlBool {
		fmt.Println("Note: Jenkins is not fully implemented yet")
		CiEnvironment.GitType = "github"

		CiEnvironment.RunnerType = "jenkins"
	} else if gitlabCiBool && gitlabCi == "true" {
		fmt.Println("Note: GitLab CI is not fully implemented yet")
	} else {
		log.Fatalln("Could not determan running environment!\nFor support please open an Issue at https://github.com/eksrvb/awesome-ci/issues")
	}

	CiEnvironment.GitInfos.DefaultBranchName = getDefaultBranch()

	return CiEnvironment
}
