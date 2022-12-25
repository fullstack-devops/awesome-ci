package main

import (
	awesomeci "awesome-ci/internal/app/awesome-ci"
)

func main() {
	//	fmt.Println("Version:\t", build.Version)
	//	fmt.Println("Commit: \t", build.CommitHash)
	//	fmt.Println("Date:   \t", build.BuildDate)
	awesomeci.AwesomeCI()
}
