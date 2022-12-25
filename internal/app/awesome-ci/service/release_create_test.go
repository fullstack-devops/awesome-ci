package service

import (
	"awesome-ci/internal/app/awesome-ci/acigithub"
	"strings"
	"testing"
)

func TestCreateFirstRelease(t *testing.T) {
	preparedReleases := &[]string{}
	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.0.1", testEnv, t) {
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{
		Version: "1.0.1",
	})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.0.1") {
		t.Errorf("Expected Relase 1.0.1 but found %s", *draftedRelease.Name)
	}
}

func TestCreateRelease_1_1_0(t *testing.T) {

	preparedReleases := &[]string{"1.0.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.1.0", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.0.1" {
		t.Errorf("Latest release should be 1.0.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.1.0") {
		t.Errorf("Expected Relase 1.1.0 but found %s", *draftedRelease.Name)
	}
}

func TestCreateRelease_1_2_0(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.1.0", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.2.0") {
		t.Errorf("Expected Relase 1.2.0 but found %s", *draftedRelease.Name)
	}
}

func TestCreateRelease_1_2_1(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.2.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.0" {
		t.Errorf("Latest release should be 1.2.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.2.1") {
		t.Errorf("Expected Relase 1.2.1 but found %s", *draftedRelease.Name)
	}
}

func TestCreateRelease_1_2_2(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.1.1", "1.1.2", "1.2.0", "1.2.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.2.2", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.2.1" {
		t.Errorf("Latest release should be 1.2.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.2.2") {
		t.Errorf("Expected Relase 1.2.2 but found %s", *draftedRelease.Name)
	}
}

func TestCreateHotfixRelease_1_1_1(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0", "1.2.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.1.1", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.0" {
		t.Errorf("Latest release should be 1.1.0 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{
		Hotfix: true,
	})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.1.1") {
		t.Errorf("Expected Relase 1.1.1 but found %s", *draftedRelease.Name)
	}
}

func TestCreateHotfixRelease_1_1_2(t *testing.T) {

	preparedReleases := &[]string{"1.0.1", "1.1.0", "1.2.0", "1.2.1", "1.1.1"}

	testEnv, cleanup := getTestEnvironment(preparedReleases, t)

	defer cleanup()

	if !resetHeadToTag("1.1.2", testEnv, t) {
		t.FailNow()
	}

	latestRelease, err := acigithub.GetLatestReleaseVersion()

	if checkError(err, t) {
		t.FailNow()
	}

	if *latestRelease.TagName != "1.1.1" {
		t.Errorf("Latest release should be 1.1.1 ... found %s", *latestRelease.TagName)
		t.FailNow()
	}

	draftedRelease := ReleaseCreate(&ReleaseCreateSet{
		Hotfix: true,
	})

	if draftedRelease == nil {
		t.Fatalf("Release Draft was not created")
	}

	defer acigithub.GithubClient.Repositories.DeleteRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, *draftedRelease.ID)

	if !*draftedRelease.Draft {
		t.Error("Release is not a Draft")
	}

	if !strings.Contains(*draftedRelease.Name, "1.1.2") {
		t.Errorf("Expected Relase 1.1.2 but found %s", *draftedRelease.Name)
	}
}
