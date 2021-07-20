package gitOnlineController

import (
	"awesome-ci/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func github_getPrNumberForBranch(branch string) int {
	url := fmt.Sprintf("%srepos/%s/pulls?state=open&head=%s", CiEnvironment.GitInfos.ApiUrl, CiEnvironment.GitInfos.FullRepo, branch[:len(branch)-1])
	respBytes := newGitHubGetRequestUnmapped(url)
	var result []models.GithubReposRepoPull

	json.Unmarshal(respBytes, &result)

	if len(result) > 0 {
		return result[0].Number
	} else {
		return 0
	}
}

func github_getLatestReleaseVersion() string {
	url := fmt.Sprintf("%srepos/%s/releases/latest", CiEnvironment.GitInfos.ApiUrl, CiEnvironment.GitInfos.FullRepo)
	result := newGitHubGetRequest(url, CiEnvironment.GitInfos.ApiToken)

	var version string
	if result["message"] == "Not Found" {
		fmt.Println("There is no release! Making initial release 0.0.0")
		version = "0.0.0"
	} else {
		version = fmt.Sprintf("%s", result["tag_name"])
	}

	return version
}

func github_createNextGitHubRelease(branch string, newReleaseVersion string, preRelease bool, uploadArtifacts string) {
	requestBody, err := json.Marshal(models.GithubNewRelease{
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

	url := fmt.Sprintf("%srepos/%s/releases", CiEnvironment.GitInfos.ApiUrl, CiEnvironment.GitInfos.FullRepo)

	respCreateRelease := newGitHubPostRequest(url, CiEnvironment.GitInfos.ApiToken, false, requestBody)
	if respCreateRelease["name"] == "Release "+newReleaseVersion {
		fmt.Println("Release " + newReleaseVersion + " sucsessfully created")
	} else {
		log.Fatalln("Somethin went worng at creating release:\n", githubErrorPrinter(respCreateRelease))
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
			respUploadArtifact := newGitHubPostRequest(newUploadUrl, CiEnvironment.GitInfos.ApiToken, true, file)
			if respUploadArtifact["name"] == releaseFileName {
				fmt.Printf("Sucsessfully uploaded asset: %s\n", releaseFileName)
			} else {
				log.Fatalln("Somethin went wrong at uploading asset:", respUploadArtifact["message"])
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

func newGitHubGetRequest(endpoint string, token string) map[string]interface{} {
	timeout := time.Duration(15 * time.Second)
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
	}

	return result
}

func newGitHubPostRequest(endpoint string, token string, isFile bool, requestBody []byte) map[string]interface{} {
	timeout := time.Duration(15 * time.Second)
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
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result
}

// run web requests to github x
func newGitHubGetRequestUnmapped(endpoint string) []byte {
	timeout := time.Duration(15 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatalln("(newGetRequest) Error at building request: ", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+CiEnvironment.GitInfos.ApiToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("(newGetRequest) Error form response:", err, resp)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["message"] == "Bad credentials" {
		log.Fatalln("Please provide the right credentials and make sure you have the right access rights!")
	}

	return b
}
