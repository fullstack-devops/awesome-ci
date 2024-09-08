package parse

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// ValidateVersion validates the given version string.
//
// Parameters:
// - version: the version string to validate.
//
// Returns:
// - err: an error if the version string is invalid.
func ValidateVersion(version string) (err error) {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	semver, err := semver.NewVersion(version)
	if err != nil {
		return
	}

	fmt.Printf("version: %s\n", version)
	fmt.Printf("  major: %d\n", semver.Major())
	fmt.Printf("  minor: %d\n", semver.Minor())
	fmt.Printf("  patch: %d\n", semver.Patch())
	if semver.Metadata() != "" {
		fmt.Printf("  metadata: %s\n", semver.Metadata())
	}
	if semver.Prerelease() != "" {
		fmt.Printf("  pre release: %s\n", semver.Prerelease())
	}

	return
}
