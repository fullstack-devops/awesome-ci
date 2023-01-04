package cmd_test

import (
	"awesome-ci/internal/app/awesome-ci/cmd"
	"log"
	"os"
	"testing"

	"github.com/spf13/cobra/doc"
)

func TestCreateCobraDocs(t *testing.T) {
	if os.Getenv("COBRA_DOCS") == "" {
		err := doc.GenMarkdownTree(cmd.RootCmd, "../../../../docs/cmd")
		if err != nil {
			log.Fatal(err)
		}
	}
}
