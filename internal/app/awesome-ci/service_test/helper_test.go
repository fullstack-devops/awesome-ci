package service_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/connect"
	scmportal "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal"
	githublocal "github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/scm-portal/github"
	"github.com/fullstack-devops/awesome-ci/internal/pkg/tools"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v56/github"
)

type TestEnvironment struct {
	ctx                 context.Context
	testOwner, testRepo string
	repoPath            string
	githubRepo          *github.Repository
	localGitRep         *git.Repository
}

func getTestEnvironment(preparedReleases *[]string, t *testing.T, ghrc *githublocal.GitHubRichClient) (testEnv *TestEnvironment, cleanup func()) {

	testRepo, ok := os.LookupEnv("ACI_TEST_REPO")

	if !ok {
		t.Errorf("env var ACI_TEST_REPO not set")
		t.FailNow()
	}

	scmLayer, err := scmportal.LoadSCMPortalLayer()
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	scmLayer.CES.
		testEnv = &TestEnvironment{
		ctx:       context.Background(),
		testOwner: strings.Split(testRepo, "/")[0],
		testRepo:  strings.Split(testRepo, "/")[1],
	}

	tempPath := path.Join(os.TempDir(), "tmp", "aci")
	if err := os.MkdirAll(tempPath, 0644); err == nil {
		testEnv.repoPath = tempPath
	} else {
		t.Errorf("Could not create tmp dir: %s", err)
		t.FailNow()
	}

	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		t.Errorf("Could not create GitHub client: %s", err)
		t.FailNow()
	}

	repo, _, err := ghrc.Client.Repositories.Get(testEnv.ctx, testEnv.testOwner, testEnv.testRepo)

	if checkError(err, t) {
		t.FailNow()
	}

	testEnv.githubRepo = repo

	gitRepo, err := git.PlainClone(testEnv.repoPath, false, &git.CloneOptions{
		URL: *repo.CloneURL,
		Auth: &http.BasicAuth{
			Username: "notneeded", // yes, this can be anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})

	if err != nil {
		t.Errorf("Could not clone Repository: %s", err)
		t.FailNow()
	}

	testEnv.localGitRep = gitRepo

	os.Chdir(testEnv.repoPath)

	cleanupReleases, err := prepareReleases(preparedReleases, testEnv, t)

	checkError(err, t)

	return testEnv, func() {
		cleanupReleases()
		os.RemoveAll(testEnv.repoPath)
	}
}

func prepareReleases(tagNames *[]string, testEnv *TestEnvironment, t *testing.T) (cleanup func(), err error) {
	//_, tagToCommitMap, err := tools.GetGitTagMaps(testEnv.localGitRep)

	tagIter, err := testEnv.localGitRep.Tags()

	if err != nil {
		return func() {}, err
	}

	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		t.Errorf("Could not create GitHub client: %s", err)
		t.FailNow()
	}

	tagIter.ForEach(func(r *plumbing.Reference) error {
		rName := r.Name().Short()

		for _, tagName := range *tagNames {
			if rName == tagName {

				commit := r.Hash().String()

				_, _, err := ghrc.Client.Repositories.CreateRelease(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, &github.RepositoryRelease{
					TagName:         &tagName,
					TargetCommitish: &commit,
					Name:            &tagName,
				})

				waitingForRelease(tagName, testEnv, t)

				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return func() {
		deleteReleases(tagNames, testEnv, t)
	}, nil
}

func waitingForRelease(tagName string, testEnv *TestEnvironment, t *testing.T) {

	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		t.Errorf("Could not create GitHub client: %s", err)
		t.FailNow()
	}

	for i := 1; i <= 10; i++ { //waiting for release becoming available

		_, _, err := ghrc.Client.Repositories.GetReleaseByTag(testEnv.ctx, testEnv.testOwner, testEnv.testRepo, tagName)
		if err != nil {
			fmt.Printf("Waiting %sms for Release to become available\n", time.Duration(i*100)*time.Millisecond)

			time.Sleep(time.Duration(i*100) * time.Millisecond)
		}
	}
}

func deleteReleases(tagNames *[]string, testConfig *TestEnvironment, t *testing.T) {

	ghrc, err := connect.ConnectToGitHub()
	if err != nil {
		t.Errorf("Could not create GitHub client: %s", err)
		t.FailNow()
	}

	for _, tagName := range *tagNames {

		release, _, err := ghrc.Client.Repositories.GetReleaseByTag(testConfig.ctx, testConfig.testOwner, testConfig.testRepo, tagName)

		if !checkError(err, t) {
			_, err = ghrc.Client.Repositories.DeleteRelease(testConfig.ctx, testConfig.testOwner, testConfig.testRepo, *release.ID)
			checkError(err, t)
		}
	}

}

func resetHeadToTag(tagName string, testEnv *TestEnvironment, t *testing.T) bool {
	_, tagToCommitMap, err := tools.GetGitTagMaps(testEnv.localGitRep)

	if checkError(err, t) {
		return false
	}
	worktree, err := testEnv.localGitRep.Worktree()

	if checkError(err, t) {
		return false
	}

	worktree.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(tagToCommitMap[tagName]),
		Mode:   git.HardReset,
	})
	return true
}

func checkError(err error, t *testing.T) bool {

	if err != nil {
		t.Error(err)
		return true
	}

	return false
}
