package rcpersist

import (
	"awesome-ci/internal/pkg/tools"
	"bufio"
	"os"
)

func checkGitIgnore() (err error) {
	if present := tools.CheckFileExists(ignoreFileName); present {
		readFile, err := os.Open(ignoreFileName)
		if err != nil {
			return err
		}

		defer readFile.Close()

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			if fileScanner.Text() == rcFileName {
				return nil
			}
		}
	}
	return ErrNotInGitIgnore
}
