package gitcontroller

type NewRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"Body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

// ReposRepoPulls match the answer for repos/___/pulls?state=all
type ReposRepoPull struct {
	Url    string `json:"url"`
	Number int    `json:"number"`
}
