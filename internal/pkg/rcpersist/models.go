package rcpersist

import "fmt"

var (
	rcTokenKey         = []byte("TZPtSIacEJG18IpqQSkTE6luYmnCNKgR")
	ErrRcFileNotExists = fmt.Errorf("rc file does not exist. please connect first")
	ErrNotInGitIgnore  = fmt.Errorf("warn: please add %s to your %s or secrets may be exposed", rcFileName, ignoreFileName)
)

const (
	rcFileName                           = ".awesomerc.yaml"
	ignoreFileName                       = ".gitignore"
	SCMPortalTypeGitHub    SCMPortalType = "github"
	SCMPortalTypeGitLab    SCMPortalType = "gitlab"
	CESTypeGitHubRunner    CESType       = "github_runner"
	CESTypeGitLabRunner    CESType       = "gitlab_runner"
	CESTypeJenkinsPipeline CESType       = "jenkins_pipeline"
	CESTypeLoMa            CESType       = "locally_manually"
)

type SCMPortalType string
type CESType string
type Token string

type RcFile struct {
	ConnectCreds  ConnectCredentials
	CESType       CESType
	SCMPortalType SCMPortalType
}

type ConnectCredentials struct {
	ServerURL  string  `yaml:"server_url"`
	Repository string  `yaml:"repository"`
	Token      Token   `yaml:"token"`
	TokenPlain *string `yaml:"tplain,omitempty"`
}
