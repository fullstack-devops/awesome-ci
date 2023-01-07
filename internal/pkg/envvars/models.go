package envvars

const EnvFile = "awesome.env"

type EnvVariables struct {
	Envs []EnvVariable
}

type EnvVariable struct {
	Name  *string
	Value *string
}
