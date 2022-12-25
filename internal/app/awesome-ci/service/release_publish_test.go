package service

import (
	"awesome-ci/internal/pkg/githubapi"
	"testing"
)

func TestPublishRelease_1_1_0(t *testing.T) {

	preparedReleases := &[]string{"1.0.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.1.0"}, testEnv, t)

	if !resetHeadToTag("1.1.0", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.0.1" {
		t.Errorf("Latest release should be 1.0.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{})

	waitingForRelease("1.1.0", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
	}
}

func TestPublishRelease_1_2_2(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0", "1.2.1", "1.1.1", "1.1.2"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.2.2"}, testEnv, t)

	if !resetHeadToTag("1.2.2", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.1" {
		t.Errorf("Latest release should be 1.2.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{})

	waitingForRelease("1.2.2", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.2" {
		t.Errorf("Latest release should be 1.2.2 ... found %s", *latestRelease.TagName)
	}
}

func TestPublishRelease_1_2_1(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.2.1"}, testEnv, t)

	if !resetHeadToTag("1.2.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.0" {
		t.Errorf("Latest release should be 1.2.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{})

	waitingForRelease("1.2.1", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.1" {
		t.Errorf("Latest release should be 1.2.1 ... found %s", *latestRelease.TagName)
	}
}

func TestPublishRelease_1_2_0(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.2.0"}, testEnv, t)

	if !resetHeadToTag("1.2.0", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{})

	waitingForRelease("1.2.0", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.0" {
		t.Errorf("Latest release should be 1.2.0 ... found %s", *latestRelease.TagName)
	}
}

func TestPublishFirstRelease(t *testing.T) {
	preparedReleases := &[]string{}
	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.0.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if latestRelease != nil || err == nil {
		t.Errorf("There should be no Relase")
	}

}

func TestPublishHotfixRelease_1_1_1(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0", "1.2.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.1.1"}, testEnv, t)

	if !resetHeadToTag("1.1.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{
		Hotfix: true,
	})

	waitingForRelease("1.1.1", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.1" {
		t.Errorf("Latest release should be 1.1.1 ... found %s", *latestRelease.TagName)
	}
}

func TestPublishHotfixRelease_1_1_2(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.1.1", "1.2.0", "1.2.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()
	defer deleteReleases(&[]string{"1.1.2"}, testEnv, t)

	if !resetHeadToTag("1.1.2", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.1" {
		t.Errorf("Latest release should be 1.1.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	ReleasePublish(&ReleasePublishSet{
		Hotfix: true,
	})

	waitingForRelease("1.1.2", testEnv, t)

	latestRelease, err = githubapi.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.2" {
		t.Errorf("Latest release should be 1.1.2 ... found %s", *latestRelease.TagName)
	}
}
