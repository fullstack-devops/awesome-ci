package cobradocs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	overrideFirstHeader *string
)

func GenerateStructuredDocs(docsPathAbs string, firstHeader string, cmd *cobra.Command) (err error) {
	if docsPathAbs == "" {
		return fmt.Errorf("docs path is empty")
	}

	if firstHeader != "" {
		overrideFirstHeader = &firstHeader
	}

	if err := os.MkdirAll(docsPathAbs, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := filepath.Walk(docsPathAbs, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to clean directory: %w", err)
	}

	err = generateDocs(cmd, docsPathAbs)
	if err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}
	return nil
}

func generateDocs(cmd *cobra.Command, parentDir string) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(parentDir, "index.md")

	// Generate the command's markdown and append it after the template
	var docContent strings.Builder
	if err := doc.GenMarkdown(cmd, &docContent); err != nil {
		return fmt.Errorf("failed to generate markdown for %s: %w", cmd.Name(), err)
	}

	content := manipulateFileContent(docContent.String())

	// Write the final content to the file
	err := os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	// Recursively generate docs for subcommands
	for _, subCmd := range cmd.Commands() {
		if !subCmd.IsAvailableCommand() || subCmd.IsAdditionalHelpTopicCommand() {
			continue
		}

		subDir := filepath.Join(parentDir, subCmd.Name())
		if err := generateDocs(subCmd, subDir); err != nil {
			return fmt.Errorf("failed to generate docs for subcommand %s: %w", subCmd.Name(), err)
		}
	}

	return nil
}

func manipulateFileContent(content string) string {
	var newLines []string
	lines := strings.Split(content, "\n")

	headline := strings.SplitAfter(lines[0], " ")
	if overrideFirstHeader != nil {
		newLines = append(newLines, "# "+*overrideFirstHeader)
		overrideFirstHeader = nil
	} else {
		newLines = append(newLines, "# "+headline[len(headline)-1])
	}

	afterSeeAlso := false

	for i, line := range lines {
		if i == 0 {
			continue
		}
		if strings.Count(line, "#") >= 2 {
			line = strings.Replace(line, "##", "#", -1)
		}
		if afterSeeAlso {
			line = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`).ReplaceAllString(line, "**$1**")
		}
		if strings.Contains(line, "SEE ALSO") {
			afterSeeAlso = true
		}

		if strings.Contains(line, "Auto generated by spf13/cobra on") {
			line = strings.Replace(line, "### Auto generated by spf13/cobra on", "##### Auto generated on", 1)
		}
		newLines = append(newLines, line)
	}

	return strings.Join(newLines, "\n")
}

func FindDirectoryUpwards(startDir, targetDir string) (string, error) {
	currentDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		targetPath := filepath.Join(currentDir, targetDir)
		if _, err := os.Stat(targetPath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", nil
}
