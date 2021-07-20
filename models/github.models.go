package models

type GithubNewRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

// ReposRepoPulls match the answer for repos/___/pulls?state=all
type GithubReposRepoPull struct {
	Url    string `json:"url"`
	Number int    `json:"number"`
}
