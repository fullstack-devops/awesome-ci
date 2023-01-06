package awesomeci

import (
	"awesome-ci/internal/app/awesome-ci/cmd"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// set to stderr becaus we want to show log information in the console and do stuff with some output
	log.SetOutput(os.Stderr)

	log.SetLevel(log.InfoLevel)
}

func AwesomeCI() {
	cmd.Execute()
}
