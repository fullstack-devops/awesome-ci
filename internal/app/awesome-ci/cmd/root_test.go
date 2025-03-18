package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/cmd"
	cobradocs "github.com/fullstack-devops/awesome-ci/internal/pkg/cobra-docs"
)

var docsPath = filepath.Join("docs", "docs", "02-CLI")

func TestRootCmd(t *testing.T) {
	if _, ok := os.LookupEnv("CI"); ok {
		t.Skip()
	}

	// create a directory to write the command information to
	rootDir, err := cobradocs.FindDirectoryUpwards(".", docsPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if rootDir == "" {
		t.Fatal("Directory not found", docsPath)
	}

	docsPath = filepath.Join(rootDir, docsPath)

	err = cobradocs.GenerateStructuredDocs(docsPath, "CLI", cmd.RootCmd)
	if err != nil {
		t.Fatal(err)
	}
}
