package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetFilesAndInfos(uploadArtifacts *string) (files []os.File, err error) {
	artifactsToUpload := strings.Split(*uploadArtifacts, ",")
	for _, artifact := range artifactsToUpload {
		var sanFilename string
		if strings.HasPrefix(artifact, "file=") {
			sanFilename = artifact[5:]
		}
		file, err := os.OpenFile(sanFilename, os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		files = append(files, *file)
	}
	return
}

func GetDefaultBranch() string {
	branch := runcmd(`git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'`, true)
	return strings.TrimSuffix(branch, "\n")
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
