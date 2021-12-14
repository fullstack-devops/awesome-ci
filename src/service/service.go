package service

import (
	"awesome-ci/src/gitController"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

var CiEnvironment gitController.CIEnvironment

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

func EvaluateEnvironment() (ciEnvironment gitController.CIEnvironment) {
	// env to check for github_runner
	githubRunnerApi, githubRunnerApiBool := os.LookupEnv("GITHUB_API_URL")
	githubRunnerRep, githubRunnerRepBool := os.LookupEnv("GITHUB_REPOSITORY")
	githubEnv, githubEnvBool := os.LookupEnv("GITHUB_ENV")
	// env to check for jenkins pipelines
	gitlabCi, gitlabCiBool := os.LookupEnv("GITLAB_CI")
	// env to check for jenkins pipelines
	_, jenkinsUrlBool := os.LookupEnv("JENKINS_URL")
	// env to check for jenkins pipelines
	// gitTypeOverride, gitTypeOverrideBool := os.LookupEnv("ACI_GIT_TYPE")

	if githubRunnerApiBool && githubRunnerRepBool && githubEnvBool {
		if !strings.HasSuffix(githubRunnerApi, "/") {
			githubRunnerApi = githubRunnerApi + "/"
		}
		ciEnvironment.GitInfos.ApiUrl = &githubRunnerApi

		ciEnvironment.GitInfos.Owner = &strings.Split(githubRunnerRep, "/")[0]
		ciEnvironment.GitInfos.Repo = &strings.Split(githubRunnerRep, "/")[1]
		githubRunnerToken, githubRunnerTokenBool := os.LookupEnv("GITHUB_TOKEN")
		if !githubRunnerTokenBool {
			log.Fatalln("Apparently you are using a GitHub-Runner.\nPlease provide the GITHUB_TOKEN!\nSee https://docs.github.com/en/actions/reference/authentication-in-a-workflow#using-the-github_token-in-a-workflow\nand https://eksrvb.github.io/awesome-ci/examples/github_actions.html")
		}
		ciEnvironment.GitInfos.ApiToken = &githubRunnerToken

		ciEnvironment.RunnerType = "github_runner"
		ciEnvironment.RunnerInfo.EnvFile = githubEnv

		gitHubTs := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubRunnerToken},
		)
		githubTc := oauth2.NewClient(context.Background(), gitHubTs)
		githubClient, err := github.NewEnterpriseClient(githubRunnerApi, githubRunnerApi, githubTc)
		if err != nil {
			log.Fatalln("error at initializing github client: ", err)
		}
		ciEnvironment.Clients.GithubClient = githubClient

	} else if jenkinsUrlBool {
		fmt.Println("Note: Jenkins is not fully implemented yet")
		ciEnvironment.RunnerType = "jenkins"

	} else if gitlabCiBool && gitlabCi == "true" {
		fmt.Println("Note: GitLab CI is not fully implemented yet")
		ciEnvironment.RunnerType = "gitlab"

		gitlabClient, err := gitlab.NewClient(os.Getenv("CI_JOB_TOKEN"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		ciEnvironment.Clients.GitlabClient = gitlabClient

	} else {
		log.Fatalln("Could not determan running environment!\nFor support please open an Issue at https://github.com/eksrvb/awesome-ci/issues")
	}

	/* if gitTypeOverrideBool {
		if strings.Contains("github gitlab", gitTypeOverride) {
			log.Printf("manual git type override requested. Using: %s", gitTypeOverride)
		} else {
			log.Fatalf("manual git type override requested. But requested type %s does not matching with github or gitlab", gitTypeOverride)
		}
	} */
	defaultBranchName := strings.Trim(getDefaultBranch(), "\n")
	ciEnvironment.GitInfos.DefaultBranchName = &defaultBranchName

	return
}

func getDefaultBranch() string {
	return runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
}
