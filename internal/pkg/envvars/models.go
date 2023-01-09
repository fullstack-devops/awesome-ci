package envvars

type EnvVariables struct {
	Envs []EnvVariable
}

type EnvVariable struct {
	Name  *string
	Value *string
}
