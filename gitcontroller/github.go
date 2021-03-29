package gitcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitConfiguration struct {
	Hostname    string
	ApiVersion  string
	Actor       string
	Repository  string
	AccessToken string
}

type NewRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"Body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

func newGetRequest(endpoint string, token string) map[string]interface{} {
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

func newPostRequest(endpoint string, token string, requestBody []byte) map[string]interface{} {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Authorization", "token "+token)

	if err == nil {
		fmt.Println(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

func github_getLatestReleaseVersion(conf GitConfiguration) string {
	url := fmt.Sprintf("%s/api/%s/repos/%s/releases/latest", conf.Hostname, conf.ApiVersion, conf.Repository)
	result := newGetRequest(url, conf.AccessToken)

	return fmt.Sprintf("%s", result["tag_name"])
}

func github_createNextGitHubRelease(conf GitConfiguration, newReleaseVersion string) {

	requestBody, err := json.Marshal(NewRelease{
		TagName:         newReleaseVersion,
		TargetCommitish: "master",
		Name:            "Release " + newReleaseVersion,
		Body:            "",
		Draft:           false,
		PreRelease:      false,
	})
	if err == nil {
		fmt.Println(err)
	}

	url := fmt.Sprintf("%s/api/%s/repos/%s/releases", conf.Hostname, conf.ApiVersion, conf.Repository)
	result := newPostRequest(url, conf.AccessToken, requestBody)
	fmt.Println(result)
}
