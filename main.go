package main

import (
	"awesome-ci-semver/gitcontroller"
	"awesome-ci-semver/semver"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func runcmd(cmd string, shell bool) string {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println(err)
			panic("some error found")
		}
		return string(out)
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func funcCreateRelease(cienv *string, overrideVersion *string, getVersionIncrease *string, isDryRun *bool) {
	var gitVersion string
	if *overrideVersion != "" {
		gitVersion = *overrideVersion
	} else {
		gitVersion = gitcontroller.GetLatestReleaseVersion(*cienv)
	}

	var patchLevel string
	if *getVersionIncrease != "" {
		patchLevel = *getVersionIncrease
	} else {
		if *cienv == "Github" {
			// Output: []string {FullString, PR, FullBranch, Orga, branch, branchBegin, restOfBranch}
			regex := `[a-zA-z ]+#([0-9]+) from (([0-9a-zA-Z-]+)/((feature|bugfix|fix)/(.+)))`
			r := regexp.MustCompile(regex)

			// mergeMessage := r.FindStringSubmatch(`Merge pull request #3 from ITC-TO-MT/feature/test-1`)
			mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
			if len(mergeMessage) > 0 {
				fmt.Printf("PR-Number: %s\n", mergeMessage[1])
				fmt.Printf("Merged branch is a %s\n", mergeMessage[5])
				patchLevel = mergeMessage[5]
			} else {
				fmt.Println("No merge message found pls make shure this regex matches: ", regex)
				fmt.Print("Example: Merge pull request #3 from some-orga/feature/awesome-feature\n\n")
				fmt.Print("If you like to set your patch level manually by flag: -level (feautre|bugfix)\n\n")
				os.Exit(1)
			}
		}
	}

	newVersion := semver.IncreaseSemVer(patchLevel, gitVersion)
	if *isDryRun {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Would writing new release: %s\n", newVersion)
	} else {
		fmt.Printf("Old version: %s\n", gitVersion)
		fmt.Printf("Writing new release: %s\n", newVersion)
		gitcontroller.CreateNextGitHubRelease(*cienv, newVersion)
	}
}

func main() {

	cienv := flag.String("cienv", "Github", "set your CI Environment for Special Featueres!\nAvalible: Jenkins, Github, Gitlab, Custom\nDefault: Github")

	createRelease := flag.NewFlagSet("createRelease", flag.ExitOnError)
	overrideVersion := createRelease.String("version", "", "override version to Update")
	getVersionIncrease := createRelease.String("level", "", "predefine version to Update")
	isDryRun := createRelease.Bool("dry-run", false, "Make dry-run before writing version to Git")

	flag.Parse()

	if os.Args[1] != "createRelease" {
		fmt.Println("expected 'createRelease', '...' or 'bar' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createRelease":
		createRelease.Parse(os.Args[2:])
		funcCreateRelease(cienv, overrideVersion, getVersionIncrease, isDryRun)
	}
}
