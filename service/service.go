package service

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

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

func getLatestCommitMessage() (infos infosMergeMessage, err error) {
	// Output: []string {FullString, PR, FullBranch, Orga, branch, branchBegin, restOfBranch}
	regex := `[a-zA-z ]+#([0-9]+) from (([0-9a-zA-Z-]+)/((feature|bugfix|fix)/(.+)))`
	r := regexp.MustCompile(regex)

	// mergeMessage := r.FindStringSubmatch(`Merge pull request #3 from test-orga/feature/test-1`)
	mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
	if len(mergeMessage) > 0 {
		infos.PRNumber = mergeMessage[1]
		infos.PatchLevel = mergeMessage[5]
		return infos, nil
	} else {
		return infos, errors.New("No merge message found pls make shure this regex matches: " + regex +
			"\nExample: Merge pull request #3 from some-orga/feature/awesome-feature" +
			"\nIf you like to set your patch level manually by flag: -level (feautre|bugfix)")
	}
}

func getDefaultBranch() string {
	return runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
}
