package ces

import "awesome-ci/internal/pkg/envvars"

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

		envFile.CloseEnvFile(ces.EnvFile)
		outFile.CloseEnvFile(ces.EnvFile)

	} else {

		for _, ev := range envVars {
			envFile.Set(ev.Name, ev.Value)
		}
		envFile.CloseEnvFile(ces.EnvFile)
	}

	return
}
