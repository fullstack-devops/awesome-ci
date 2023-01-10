package ces

import (
	"awesome-ci/internal/pkg/envvars"

	log "github.com/sirupsen/logrus"
)

func (ces CES) ExportAsEnv(envVars []KeyValue) (err error) {
	envFile, err := envvars.OpenEnvFile(ces.EnvFile)
	if err != nil {
		return
	}

	if ces.OutFile != nil {
		outFile, err := envvars.OpenEnvFile(ces.EnvFile)
		if err != nil {
			return err
		}
		for _, ev := range envVars {
			envFile.Set(ev.Name, ev.Value)
			outFile.Set(ev.Name, ev.Value)
		}

		if err := envFile.CloseEnvFile(ces.EnvFile); err != nil {
			log.Warnln(err)
			return err
		}
		log.Tracef("succfully written env variables to %s", ces.EnvFile)

		if err := outFile.CloseEnvFile(ces.EnvFile); err != nil {
			log.Warnln(err)
			return err
		}
		log.Tracef("succfully written env variables to %s", *ces.OutFile)

	} else {

		for _, ev := range envVars {
			envFile.Set(ev.Name, ev.Value)
		}
		envFile.CloseEnvFile(ces.EnvFile)
	}

	return
}
