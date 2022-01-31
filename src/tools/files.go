package tools

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

type UploadArtifact struct {
	File os.File
	Name string
}

type UploadAsset struct {
	File  os.File
	Name  string
	Infos fs.FileInfo
	IsZip bool
	Hash  []byte
}

// deprecated
func GetFilesAndInfos(uploadArtifacts *string) (artifacts []UploadArtifact, err error) {
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
		info, _ := file.Stat()
		artifacts = append(artifacts, UploadArtifact{
			File: *file,
			Name: info.Name(),
		})
	}
	return
}

func GetAsstes(uploadArtifacts *string, check bool) (asset []UploadAsset, err error) {
	artifactsToUpload := strings.Split(*uploadArtifacts, ",")
	for _, artifact := range artifactsToUpload {
		assetName := strings.Split(artifact, "=")

		switch assetName[0] {
		case "file":
			file, err := os.OpenFile(assetName[1], os.O_RDWR, 0644)
			if err != nil {
				return nil, err
			}
			info, err := file.Stat()
			if err != nil {
				return nil, err
			}
			// if check, then close file, else return file interface for upload
			if check {
				file.Close()
			} else {
				hash, _ := getHashFromFile(assetName[1])
				asset = append(asset, UploadAsset{
					File:  *file,
					Name:  info.Name(),
					Infos: info,
					IsZip: false,
					Hash:  hash,
				})
			}
		case "files":
			// todo
		case "zip":
			// todo
		default:
			return nil, fmt.Errorf("not an valid asset format")
		}
	}

	return
}

func getHashFromFile(name string) (hash []byte, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}

	return h.Sum(nil), nil
}

func CheckIsFile(name string) (body string, err error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
