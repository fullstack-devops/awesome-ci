package models

type CIEnvironment struct {
	GitType  string
	GitInfos struct {
		ApiUrl            string
		ApiToken          string
		Repo              string
		Orga              string
		IsOrg             bool
		FullRepo          string
		DefaultBranchName string
	}
	RunnerType string
	RunnerInfo struct {
		EnvFile string
	}
}

type BuildEnvironmentVariable struct {
	Key   string
	Value string
}

type BuildInfos struct {
	PrNumber    int
	Version     string
	NextVersion string
	PatchLevel  string
}
