package detect

type CiType string

const (
	NormalEnv           CiType = "NormalEnv"
	GitHubActionsRunner CiType = "GitHubActionsRunner"
)

type EnvVariables struct {
	Envs   []EnvVariable
	CiType CiType
}

type EnvVariable struct {
	Name  *string
	Value *string
}

type ConnectCredentials struct {
	ServerUrl  string
	Repository string
	Token      string
}
