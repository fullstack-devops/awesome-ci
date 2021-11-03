package semver

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IncreaseSemVer increases a given version by naming, see docs
func IncreaseVersion(versionNaming string, version string) (incresedVersion string, err error) {
	incresedVersion = version

	if strings.HasPrefix(versionNaming, "major") {
		incresedVersion = major(version)
	} else if strings.HasPrefix(versionNaming, "feature") {
		incresedVersion = minor(version)
	} else if strings.HasPrefix(versionNaming, "bugfix") {
		incresedVersion = patch(version)
	} else {
		err = errors.New("Could not determan witch version to set. Given first string does'n start with release, feature or bugfix")
	}

	return
}

func major(version string) string {
	var newVersion string
	splitedVersion := strings.Split(version, ".")

	major, err := strconv.Atoi(splitedVersion[0])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	newMajor := (major + 1)

	newVersion = fmt.Sprintf("%d.0.0", newMajor)
	return newVersion
}
func minor(version string) string {
	var newVersion string
	splitedVersion := strings.Split(version, ".")

	minor, err := strconv.Atoi(splitedVersion[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	newMinor := (minor + 1)

	newVersion = fmt.Sprintf("%s.%d.0", splitedVersion[0], newMinor)
	return newVersion
}
func patch(version string) string {
	var newVersion string
	splitedVersion := strings.Split(version, ".")

	patch, err := strconv.Atoi(splitedVersion[2])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	newPatch := (patch + 1)

	newVersion = fmt.Sprintf("%s.%s.%d", splitedVersion[0], splitedVersion[1], newPatch)
	return newVersion
}
