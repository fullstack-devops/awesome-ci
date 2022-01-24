package semver

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// IncreaseSemVer increases a given version by naming, see docs
func IncreaseVersion(patchLevel string, version string) (incresedVersion string, err error) {
	incresedVersion = version

	switch patchLevel {
	case "fix", "bugfix", "dependabot":
		incresedVersion, err = patch(version)
	case "feature", "feat":
		incresedVersion, err = minor(version)
	case "major":
		incresedVersion, err = major(version)
	default:
		incresedVersion, err = patch(version)
		log.Println("could not determan witch version to set. given first string does'n start with release, feature or bugfix")
		log.Println("used minimal patch version!")
	}
	return
}

func major(version string) (newVersion string, err error) {
	splitedVersion := strings.Split(version, ".")

	major, err := strconv.Atoi(splitedVersion[0])
	if err != nil {
		return
	}
	newMajor := (major + 1)

	newVersion = fmt.Sprintf("%d.0.0", newMajor)
	return
}
func minor(version string) (newVersion string, err error) {
	splitedVersion := strings.Split(version, ".")

	minor, err := strconv.Atoi(splitedVersion[1])
	if err != nil {
		return
	}
	newMinor := (minor + 1)

	newVersion = fmt.Sprintf("%s.%d.0", splitedVersion[0], newMinor)
	return
}
func patch(version string) (newVersion string, err error) {
	splitedVersion := strings.Split(version, ".")

	patch, err := strconv.Atoi(splitedVersion[2])
	if err != nil {
		return
	}
	newPatch := (patch + 1)

	newVersion = fmt.Sprintf("%s.%s.%d", splitedVersion[0], splitedVersion[1], newPatch)
	return
}
