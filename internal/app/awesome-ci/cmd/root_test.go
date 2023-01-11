package cmd_test

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fullstack-devops/awesome-ci/internal/app/awesome-ci/cmd"

	"github.com/spf13/cobra/doc"
)

var docsPath = "../../../../docs/cmd"

const fmTemplate = `---
layout: default
nav_order: 2
parent: CLI
date: %s
title: "%s"
---
`

func TestCreateCobraDocs(t *testing.T) {
	/* if os.Getenv("COBRA_DOCS") == "" {
		err := doc.GenMarkdownTree(cmd.RootCmd, "../../../../docs/cmd")
		if err != nil {
			log.Fatal(err)
		}
	} */
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, docsPath, filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func filePrepender(filename string) string {
	now := time.Now().Format(time.RFC3339)
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))

	if base != "awesome-ci" {
		base = strings.Replace(base, "awesome-ci", "", -1)
	}
	return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1))
}

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	return "/commands/" + strings.ToLower(base) + "/"
}
