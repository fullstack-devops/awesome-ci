package semver_test

import (
	"testing"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/semver"
)

func TestPatchLevels(t *testing.T) {
	testPatchLevelFormBranch(t, semver.Bugfix, "bugfix/test-#1")
	testPatchLevelFormBranch(t, semver.Bugfix, "something/test-#1")
	testPatchLevelFormBranch(t, semver.Bugfix, "dependabot/test-#1")
	testPatchLevelFormBranch(t, semver.Feature, "feature/test-#1")
	testPatchLevelFormBranch(t, semver.Feature, "feat/test-#1")
	testPatchLevelFormBranch(t, semver.Major, "major/test-#1")
}

func TestVersionIncrease(t *testing.T) {
	testVersionIncrease(t, semver.Bugfix, "0.0.1", "0.0.2")
	testVersionIncrease(t, semver.Bugfix, "0.2.7", "0.2.8")
	testVersionIncrease(t, semver.Bugfix, "6.5.1", "6.5.2")

	testVersionIncrease(t, semver.Feature, "6.5.1", "6.6.0")

	testVersionIncrease(t, semver.Major, "6.5.1", "7.0.0")
}

func TestInvalidBranchName(t *testing.T) {
	_, err := semver.ParsePatchLevelFormBranch("something-else")
	if err == nil {
		t.Errorf("shout fail at wrong branch name")
		t.FailNow()
	}
}

func testPatchLevelFormBranch(t *testing.T, patchLevelExpected semver.PatchLevel, branchName string) {
	final, err := semver.ParsePatchLevelFormBranch(branchName)
	if err != nil && err != semver.ErrUseMinimalPatchVersion {
		t.Errorf("%v", err)
		t.FailNow()
	}
	if final != patchLevelExpected {
		t.Errorf("got not expected PatchLevel, should be %s is: %s", patchLevelExpected, final)
		t.FailNow()
	}
}

func testVersionIncrease(t *testing.T, patchLevel semver.PatchLevel, oldVersion string, versionExpected string) {
	final, err := semver.IncreaseVersion(patchLevel, oldVersion)
	if err != nil && err != semver.ErrUseMinimalPatchVersion {
		t.Errorf("%v", err)
		t.FailNow()
	}
	if final != versionExpected {
		t.Errorf("got not expected Version, should be %s is: %s", versionExpected, final)
		t.FailNow()
	}
}
