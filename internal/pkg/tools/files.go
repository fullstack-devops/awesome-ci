package tools

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrEmptyString        = fmt.Errorf("empty string")
	ErrExpextedFileNotDir = fmt.Errorf("expected file and not a directory")
)

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
