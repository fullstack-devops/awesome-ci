package service

import (
	"awesome-ci/gitcontroller"
	"awesome-ci/semver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func CreateRelease(cienv *string, overrideVersion *string, getVersionIncrease *string, isDryRun *bool, publishNpm *string) {
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

	if *publishNpm != "" {
		// check if subfolder has slash
		pathToSource := *publishNpm
		if !strings.HasSuffix(*publishNpm, "/") {
			pathToSource = *publishNpm + "/"
		}
		fmt.Printf("Puplishing npm packages under path: %s\n", pathToSource)
		npmPublish(pathToSource, newVersion)
	}
}

func npmPublish(pathToSource string, newVersion string) {

	// opening package.json
	jsonFile, err := os.Open(pathToSource + "package.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer jsonFile.Close()

	var result map[string]interface{}
	json.NewDecoder(jsonFile).Decode(&result)

	result["version"] = newVersion

	b, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// writing result to package.json
	err = ioutil.WriteFile(pathToSource+"package.json", b, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	command := exec.Command("npm", "publish", pathToSource, "--tag", fmt.Sprintf("%s@%s", result["name"], newVersion))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// Run the command
	command.Run()
}
