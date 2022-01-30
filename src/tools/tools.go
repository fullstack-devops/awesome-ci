package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetDefaultBranch() string {
	branch := runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
	return strings.TrimSuffix(branch, "\n")
}

func DevideOwnerAndRepo(fullRepo string) (owner string, repo string) {
	owner = strings.ToLower(strings.Split(fullRepo, "/")[0])
	repo = strings.ToLower(strings.Split(fullRepo, "/")[1])
	return
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
