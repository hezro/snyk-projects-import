package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const base_url = "https://api.snyk.io/api/v1"

// Define flags for the utility.
var (
	snyk_token         = kingpin.Flag("token", "Snyk API Token").Required().String()
	git_integration_id = kingpin.Flag("gitId", "Git/SCM integration ID").Required().String()
	org_id             = kingpin.Flag("orgId", "Snyk Target Organization ID").Required().String()
	repo_owner         = kingpin.Flag("owner", "Account owner of the repository").Required().String()
	repo_name          = kingpin.Flag("repoName", "Name of the Repo").Required().String()
	branch_name        = kingpin.Flag("branchName", "Name of the Branch").Required().String()
	file_path          = kingpin.Flag("filePath", "Relative path to one or more files").String()
)

// Define the structs for the json fields.
type DataAttributes struct {
	Target Target  `json:"target"`
	Files  []Files `json:"files"`
}
type Target struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}
type Files struct {
	Path string `json:"path"`
}

// Set up the http client with a 10 second timeout.
func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

// Import file(s) or the whole repo into Snyk.
func snyk_import(client *http.Client, snyk_token string, git_integration_id string, org_id string, repo_owner string, repo_name string, branch_name string, file_path []string) []byte {
	url := base_url + "/org/" + org_id + "/integrations/" + git_integration_id + "/import"

	// Build the json data.
	target := DataAttributes{}
	target.Target.Owner = repo_owner
	target.Target.Name = repo_name
	target.Target.Branch = branch_name

	// If files are passed in then add them to the json data.
	file_count := len(file_path)
	if file_count > 0 {
		data := Files{}
		for i := range file_path {
			data = (Files{file_path[i]})
			target.Files = append(target.Files, data)
		}
	}

	json_data, _ := json.Marshal(target)
	body := bytes.NewBuffer(json_data)

	// Send the POST request
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", snyk_token)

	if err != nil {
		log.Fatalf("| Error | Failed: %+v", err)
	}

	// Send the POST request
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("| Error | Failed sending request to API endpoint. |  %+v \n", err)
	}
	defer response.Body.Close()

	response_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("| Error | Could not parse response body. |  %+v \n", err)
	}

	status_code := response.StatusCode
	if status_code != 201 {
		log.Printf("| Error |  Import failed to POST | Status Code: %s \n", strconv.Itoa((status_code)))
	}

	return response_body
}

// Create a files slice
func build_file_slice(files string) []string {
	slice := strings.Split(files, ",")
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}
	return slice
}

func main() {
	// Parse command line arguments
	kingpin.Parse()

	// Only create the files slice if we have any file(s) to import
	var files []string
	if len(*file_path) > 0 {
		files = build_file_slice(*file_path)
	}

	// Create the http client
	client := httpClient()

	// Start the import
	import_projects := snyk_import(client, *snyk_token, *git_integration_id, *org_id, *repo_owner, *repo_name, *branch_name, files)

	// Report if the POST import was successful
	body_len := len(import_projects)
	if body_len == 2 {
		log.Printf("| Success | POST request submitted successfully items from %s to orgId: %s\n", *repo_name, *org_id)
	} else {
		log.Printf("| Error |  Import failed to POST | Body: %s \n", import_projects)
	}
}
