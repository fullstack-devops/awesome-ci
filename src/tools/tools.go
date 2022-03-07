package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func GetGitTagMaps(gitRepo *git.Repository) (commitToTagMap map[string][]string, tagToCommitMap map[string]string, err error) {
	tagToCommitMap = make(map[string]string)
	commitToTagMap = make(map[string][]string)

	tags, err := gitRepo.Tags()

	if err != nil {
		return nil, nil, err
	}

	tags.ForEach(func(r *plumbing.Reference) error {

		tagList, exists := commitToTagMap[r.Hash().String()]

		if !exists {
			tagList = make([]string, 0)
			commitToTagMap[r.Hash().String()] = tagList
		}

		commitToTagMap[r.Hash().String()] = append(tagList, r.Name().Short())
		tagToCommitMap[r.Name().Short()] = r.Hash().String()
		fmt.Printf("Tag %s %s\n", r.Hash().String(), r.Name().Short())
		return nil
	})

	return commitToTagMap, tagToCommitMap, nil
}
