package connect

import "github.com/google/go-github/v49/github"

type ServerType string
type Token string

type RcFile struct {
	Type         ServerType
	ConnectCreds ConnectCredentials
}

type GitHubRichClient struct {
	Client     *github.Client
	Owner      string
	Repository string
}

type ConnectCredentials struct {
	ServerUrl  string
	Repository string
	Token      Token
	TokenPlain string `yaml:"omitempty"`
}
