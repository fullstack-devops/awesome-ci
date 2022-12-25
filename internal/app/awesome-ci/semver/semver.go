package semver

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type PatchLevel string

const (
	Bugfix  PatchLevel = "bugfix"
	Feature PatchLevel = "feature"
	Major   PatchLevel = "major"
)

// IncreaseSemVer increases a given version by naming, see docs
func IncreaseVersion(patchLevel PatchLevel, version string) (incresedVersion string, err error) {
	incresedVersion = version

	switch patchLevel {
	case Bugfix:
		incresedVersion, err = patch(version)
	case Feature:
		incresedVersion, err = minor(version)
	case Major:
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

func ParsePatchLevel(branchName string) PatchLevel {
	patchLevel := "bugfix"

	if strings.Index(branchName, "/") > 0 {
		patchLevel = branchName[:strings.Index(branchName, "/")]
	}

	switch patchLevel {
	case "bugfix", "fix", "dependabot":
		return Bugfix
	case "feature", "feat":
		return Feature
	case "major":
		return Major
	default:
		log.Println("could not determan witch version to set. given first string does'n start with major, feature or bugfix")
		log.Println("using minimal patch version!")
		return Bugfix
	}

}
