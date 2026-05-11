package relctl

import (
	"os"

	"github.com/layer87-labs/relctl/internal/app/relctl/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// set to stderr becaus we want to show log information in the console and do stuff with some output
	log.SetOutput(os.Stderr)

	log.SetLevel(log.InfoLevel)
}

func Execute() {
	cmd.Execute()
}
