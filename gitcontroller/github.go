package gitcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitConfiguration struct {
	ApiUrl            string
	Repository        string
	AccessToken       string
	DefaultBranchName string
}

type NewRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"Body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

func newGitHubGetRequest(endpoint string, token string) map[string]interface{} {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+token)
	if err != nil {
		log.Fatalln("(newGetRequest) Error at building request: ", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("(newGetRequest) Error form response:", err, resp)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["message"] == "Bad credentials" {
		log.Fatalln("Please provide the right credentials and make sure you have the right access rights!")
		os.Exit(1)
	}

	return result
}

func newGitHubPostRequest(endpoint string, token string, isFile bool, requestBody []byte) map[string]interface{} {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln("(newPostRequest) Error at building request: ", err)
	}
	if isFile {
		request.Header.Set("Content-Type", "application/octet-stream")
	} else {
		request.Header.Set("Accept", "application/vnd.github.v3+json")
	}
	request.Header.Set("Authorization", "token "+token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln("(newPostRequest) Error form response:", err, resp)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

func github_getLatestReleaseVersion(apiUrl string, repository string, accessToken string) string {
	if !strings.HasSuffix(apiUrl, "/") {
		apiUrl = apiUrl + "/"
	}
	url := fmt.Sprintf("%srepos/%s/releases/latest", apiUrl, repository)
	result := newGitHubGetRequest(url, accessToken)

	var version string
	if result["message"] == "Not Found" {
		fmt.Println("There is no release! Making initial release 0.0.0")
		version = "0.0.0"
	} else {
		version = fmt.Sprintf("%s", result["tag_name"])
	}

	return version
}

func github_createNextGitHubRelease(apiUrl string, repository string, accessToken string, branch string, newReleaseVersion string, preRelease bool, uploadArtifacts string) {
	requestBody, err := json.Marshal(NewRelease{
		TagName:         newReleaseVersion,
		TargetCommitish: strings.Trim(branch, "\n"),
		Name:            "Release " + newReleaseVersion,
		Body:            "",
		Draft:           false,
		PreRelease:      preRelease,
	})
	if err != nil {
		fmt.Println("(github_createNextGitHubRelease) Error building requestBody: ", err)
	}

	url := fmt.Sprintf("%s/repos/%s/releases", apiUrl, repository)
	log.Println("url for creating release:", url)
	respCreateRelease := newGitHubPostRequest(url, accessToken, false, requestBody)
	if respCreateRelease["name"] == "Release "+newReleaseVersion {
		fmt.Println("Release " + newReleaseVersion + " sucsessfully created")
	} else {
		log.Fatalln("Somethin went worng at creating release:\n", githubErrorPrinter(respCreateRelease))
		os.Exit(1)
	}

	if uploadArtifacts != "" {
		log.Printf("Uploading artifacts from: %s\n", uploadArtifacts)

		artifactsToUpload := strings.Split(uploadArtifacts, ",")

		for _, artifact := range artifactsToUpload {
			file, err := ioutil.ReadFile(artifact)
			if err != nil {
				log.Fatal(err)
			}

			releaseFileName := artifact[strings.LastIndex(artifact, "/")+1:]

			uploadUrl := fmt.Sprintf("%s", respCreateRelease["upload_url"])
			newUploadUrl := strings.Replace(uploadUrl, "{?name,label}", "?name="+releaseFileName, -1)
			// log.Println("url for uploading asset to release:", newUploadUrl)
			respUploadArtifact := newGitHubPostRequest(newUploadUrl, accessToken, true, file)
			if respUploadArtifact["name"] == releaseFileName {
				fmt.Printf("Sucsessfully uploaded asset: %s\n", releaseFileName)
			} else {
				log.Fatalln("Somethin went wrong at uploading asset:", respUploadArtifact["message"])
				os.Exit(1)
			}
		}
	}
}

func githubErrorPrinter(responseErrors map[string]interface{}) string {
	var errors []map[string]interface{}
	outputString := fmt.Sprintln(responseErrors["message"])

	b, err := json.Marshal(responseErrors["errors"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &errors)
	for index := range errors {
		if errors[index]["code"] == "custom" {
			outputString = outputString + fmt.Sprintf("code: %s => message: %s\n", errors[index]["code"], errors[index]["message"])
		} else {
			outputString = outputString + fmt.Sprintf("code: %s => message: %s\n", errors[index]["code"], errors[index])
		}
	}
	return outputString
}
