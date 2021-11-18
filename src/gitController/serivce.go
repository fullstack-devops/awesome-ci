package gitController

import (
	"os"
	"strings"
)

func getFilesAndInfos(uploadArtifacts *string) (files []os.File, err error) {
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
