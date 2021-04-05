package gitcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type GitConfiguration struct {
	ApiUrl      string
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
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+token)
	if err != nil {
		fmt.Println("(newGetRequest) Error at building request: ", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("(newGetRequest) Error form response:", err, resp)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["message"] == "Bad credentials" {
		fmt.Println("Please provide the right credentials and make sure you have the right access rights!")
		os.Exit(1)
	}

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
	if err != nil {
		fmt.Println("(newPostRequest) Error at building request: ", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("(newPostRequest) Error form response:", err, resp)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

func github_getLatestReleaseVersion(conf GitConfiguration) string {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", conf.ApiUrl, conf.Repository)
	result := newGetRequest(url, conf.AccessToken)

	var version string
	if result["message"] == "Not Found" {
		fmt.Println("There is no release! Making initial release 0.0.0")
		version = "0.0.0"
	} else {
		version = fmt.Sprintf("%s", result["tag_name"])
	}

	return version
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
	if err != nil {
		fmt.Println("(github_createNextGitHubRelease) Error building requestBody", err)
	}

	url := fmt.Sprintf("%s/repos/%s/releases", conf.ApiUrl, conf.Repository)
	result := newPostRequest(url, conf.AccessToken, requestBody)

	if result["name"] == "Release "+newReleaseVersion {
		fmt.Println("Release " + newReleaseVersion + " sucsessfully created")
	} else {
		fmt.Println("Somethin went worng at creating release!")
	}
}
