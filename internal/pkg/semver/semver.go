package semver

import (
	"fmt"
	"strconv"
	"strings"
)

type PatchLevel string

const (
	Bugfix  PatchLevel = "bugfix"
	Feature PatchLevel = "minor"
	Major   PatchLevel = "major"
)

var (
	ErrUseMinimalPatchVersion = fmt.Errorf("could not determan witch version to set. given first string does'n start with major, feature or bugfix -> using minimal patch version")
	ErrInvalidBranchName      = fmt.Errorf("the branch name does not match the prescribed spelling")
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

func ParsePatchLevelFormBranch(branchName string) (patchLevel PatchLevel, err error) {
	patchLevel = "bugfix"

	if strings.Index(branchName, "/") > 0 {
		return ParsePatchLevel(branchName[:strings.Index(branchName, "/")])
	} else {
		return Bugfix, ErrInvalidBranchName
	}
}

func ParsePatchLevel(alias string) (PatchLevel, error) {
	switch alias {

	case "patch", "bugfix", "fix", "dependabot":
		return Bugfix, nil

	case "minor", "feature", "feat":
		return Feature, nil

	case "major":
		return Major, nil

	default:
		return Bugfix, ErrUseMinimalPatchVersion
	}
}
