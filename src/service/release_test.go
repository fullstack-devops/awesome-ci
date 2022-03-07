package service

import (
	"awesome-ci/src/acigithub"
	"testing"
)

func TestCreateRelease_1_1_0(t *testing.T) {

	preparedReleases := &[]string{"1.0.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.1.0"}, testEnv, t)

	if !resetHeadToTag("1.1.0", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion(testEnv.testOwner, testEnv.testRepo)

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.0.1" {
		t.Errorf("Latest release should be 1.0.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{})

	waitingForRelease("1.1.0", testEnv, t)

	latestRelease, err = acigithub.GetLatestReleaseVersion(testEnv.testOwner, testEnv.testRepo)

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
	}
}

func TestFirstRelease(t *testing.T) {
	preparedReleases := &[]string{}
	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.0.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion(testEnv.testOwner, testEnv.testRepo)

	if latestRelease != nil || err == nil {
		t.Errorf("There should be no Relase")
	}

}
