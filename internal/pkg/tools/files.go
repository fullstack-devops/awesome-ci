package tools

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

var (
	ErrEmptyString        = fmt.Errorf("empty string")
	ErrExpextedFileNotDir = fmt.Errorf("expected file and not a directory")
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

// ReadFileToString can find a given file and returns it as a string
// if the input is not an empty string, the given input is returned
func ReadFileToString(fileOrString string) (result string, err error) {
	if fileOrString != "" {

		fileInfo, err := os.Stat(fileOrString)
		if err != nil && err != os.ErrNotExist {
			return "", err
		} else if err == os.ErrNotExist {
			return fileOrString, nil
		}

		if fileInfo.IsDir() {
			return "", ErrExpextedFileNotDir
		}

		if file, err := os.ReadFile(fileOrString); err != nil {
			return "", err
		} else {
			return string(file), err
		}

	} else {
		err = ErrEmptyString
	}

	return
}

func GetAsset(assetLocation string) (asset *UploadAsset, err error) {
	assetName := strings.Split(assetLocation, "=")

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
		hash, _ := getHashFromFile(assetName[1])

		return &UploadAsset{
			File:  *file,
			Name:  info.Name(),
			Infos: info,
			IsZip: false,
			Hash:  hash,
		}, nil

	case "files":
		// todo
	case "zip":
		// todo
	default:
		return nil, fmt.Errorf("not an valid asset format")
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

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
