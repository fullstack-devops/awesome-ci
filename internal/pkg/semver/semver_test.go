package semver_test

import (
	"testing"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"
	"github.com/stretchr/testify/require"
)

func TestPatchLevels(t *testing.T) {
	patchLevel, err := semver.ParsePatchLevelFormBranch("bugfix/test-#1")
	require.NoError(t, err)
	require.Equal(t, semver.Bugfix, patchLevel)

	_, err = semver.ParsePatchLevelFormBranch("something/test-#1")
	require.Error(t, err)
	require.Equal(t, semver.ErrUseMinimalPatchVersion, err)

	patchLevel, err = semver.ParsePatchLevelFormBranch("dependabot/test-#1")
	require.NoError(t, err)
	require.Equal(t, semver.Bugfix, patchLevel)

	patchLevel, err = semver.ParsePatchLevelFormBranch("feature/test-#1")
	require.NoError(t, err)
	require.Equal(t, semver.Feature, patchLevel)

	patchLevel, err = semver.ParsePatchLevelFormBranch("feat/test-#1")
	require.NoError(t, err)
	require.Equal(t, semver.Feature, patchLevel)

	patchLevel, err = semver.ParsePatchLevelFormBranch("major/test-#1")
	require.NoError(t, err)
	require.Equal(t, semver.Major, patchLevel)
}

func TestVersionIncrease(t *testing.T) {
	nextVersion, err := semver.IncreaseVersion(semver.Bugfix, "0.0.1")
	require.NoError(t, err)
	require.Equal(t, "0.0.2", nextVersion)

	nextVersion, err = semver.IncreaseVersion(semver.Bugfix, "0.2.7")
	require.NoError(t, err)
	require.Equal(t, "0.2.8", nextVersion)

	nextVersion, err = semver.IncreaseVersion(semver.Bugfix, "6.5.1")
	require.NoError(t, err)
	require.Equal(t, "6.5.2", nextVersion)

	nextVersion, err = semver.IncreaseVersion(semver.Feature, "6.5.1")
	require.NoError(t, err)
	require.Equal(t, "6.6.0", nextVersion)

	nextVersion, err = semver.IncreaseVersion(semver.Major, "6.5.1")
	require.NoError(t, err)
	require.Equal(t, "7.0.0", nextVersion)
}

func TestInvalidBranchName(t *testing.T) {
	_, err := semver.ParsePatchLevelFormBranch("something-else")
	require.Error(t, err)
}
