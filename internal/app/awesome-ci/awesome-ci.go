package awesomeci

import (
	"awesome-ci/internal/app/awesome-ci/cmd"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func AwesomeCI() {
	cmd.Execute()
}
