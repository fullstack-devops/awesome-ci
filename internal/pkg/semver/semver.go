package semver

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
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

// IncreaseVersion increases a given version by the specified patch level.
//
// Parameters:
// - patchLevel: the patch level to increase the version by (Bugfix, Feature, Major).
// - version: the current version string.
//
// Returns:
// - incresedVersion: the increased version string.
// - err: an error if the patch level is invalid or if there is an error increasing the version.
func IncreaseVersion(patchLevel PatchLevel, version string) (incresedVersion string, err error) {
	incresedVersion = version

	semVer, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	switch patchLevel {
	case Bugfix:
		return semVer.IncPatch().String(), nil
	case Feature:
		return semVer.IncMinor().String(), nil
	case Major:
		return semVer.IncMajor().String(), nil
	}

	incresedVersion = semVer.String()
	return
}

func ParsePatchLevelFormBranch(branchName string) (patchLevel PatchLevel, err error) {
	patchLevel = "bugfix" // default patch level

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
