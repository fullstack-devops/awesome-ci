package service

import (
	"awesome-ci/gitcontroller"
	"awesome-ci/semver"
	"log"
	"os"
	"regexp"
)

func GetBuildInfos(cienv string, overrideVersion *string, patchLevel *string, output *string) {
	var gitVersion string
	if *overrideVersion != "" {
		gitVersion = *overrideVersion
	} else {
		gitVersion = gitcontroller.GetLatestReleaseVersion(cienv)
	}

	if *patchLevel != "" {
		log.Println("patchLevel has Override:", *patchLevel)
	} else {
		if cienv == "Github" {
			// Output: []string {FullString, PR, FullBranch, Orga, branch, branchBegin, restOfBranch}
			regex := `[a-zA-z ]+#([0-9]+) from (([0-9a-zA-Z-]+)/((feature|bugfix|fix)/(.+)))`
			r := regexp.MustCompile(regex)

			// mergeMessage := r.FindStringSubmatch(`Merge pull request #3 from ITC-TO-MT/feature/test-1`)
			mergeMessage := r.FindStringSubmatch(runcmd(`git log -1 --pretty=format:"%s"`, true))
			if len(mergeMessage) > 0 {
				log.Printf("PR-Number: %s\n", mergeMessage[1])
				log.Printf("Merged branch is a %s\n", mergeMessage[5])
				patchLevel = &mergeMessage[5]
			} else {
				log.Println("No merge message found pls make shure this regex matches: ", regex)
				log.Print("Example: Merge pull request #3 from some-orga/feature/awesome-feature\n\n")
				log.Print("If you like to set your patch level manually by flag: -level (feautre|bugfix)\n\n")
				os.Exit(1)
			}
		}
	}
	log.Printf("New Version: %s", semver.IncreaseSemVer(*patchLevel, gitVersion))
}
