package main

import (
	"awesom-ci-semver/semver"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

/*
func semver(version string, mmb string) string {
	var newVersion string

	switch mmb {
	case "major":
		newVersion = version
	}
	return newVersion
} //git describe --tags `git rev-list --tags --max-count=1` */

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

func main() {
	getVersion := flag.NewFlagSet("getver", flag.ExitOnError)
	// overrideVersion := getVersion.String("oldVersion", "", "predefine version to Update")

	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "getver":
		getVersion.Parse(os.Args[2:])

		gitVersion := runcmd("git describe --tags $(git rev-list --tags --max-count=1)", true)

		fmt.Printf("%s", semver.Major(gitVersion))
	}
	/*
		textPtr := flag.String("getver", "", "Gives you the next version with semver")
		metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};.")
		uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
		flag.Parse() */

	// fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *textPtr, *metricPtr, *uniquePtr)
}
