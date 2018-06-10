package in

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jghiloni/spring-initializr-resource"
)

// Command will perform the in operation and download the appropriate artifacts for a given version
type Command struct {
	Client *http.Client
}

var emptyResponse = Response{}

// Run is the main unit of work for the Command
func (command *Command) Run(destinationDir string, request Request) (Response, error) {
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return emptyResponse, err
	}

	queryParams := url.Values{}
	setValueOrDefault(&queryParams, "type", request.Params.Type, "maven-project")
	setValue(&queryParams, "packaging", request.Params.Packaging)
	setValue(&queryParams, "language", request.Params.Language)
	setValue(&queryParams, "dependencies", request.Params.Dependencies)
	setValue(&queryParams, "javaVersion", request.Params.JDKVersion)
	setValue(&queryParams, "bootVersion", request.Version.ID)
	setValue(&queryParams, "groupId", request.Params.GroupID)
	setValue(&queryParams, "artifactId", request.Params.ArtifactID)
	setValue(&queryParams, "version", request.Params.Version)
	setValue(&queryParams, "name", request.Params.Name)
	setValue(&queryParams, "description", request.Params.Description)
	setValue(&queryParams, "packageName", request.Params.PackageName)

	targetURL := request.Source.URL

	endpoint := ""
	switch request.Params.Type {
	case "maven-project", "gradle-project":
		endpoint = "/starter.zip"
	case "maven-build":
		endpoint = "/pom.xml"
	case "gradle-build":
		endpoint = "/build.gradle"
	}

	targetURL.Path = path.Join(targetURL.Path, endpoint)
	targetURL.RawQuery = queryParams.Encode()

	httpRequest, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		return emptyResponse, err
	}

	httpRequest.Header.Add("Accept", initializr.AcceptHeader)
	httpResponse, err := command.Client.Do(httpRequest)
	if err != nil {
		return emptyResponse, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return emptyResponse, err
	}

	if httpResponse.StatusCode != 200 {
		return emptyResponse, fmt.Errorf("Expected 200 OK, got %d %s with message %s", httpResponse.StatusCode, httpResponse.Status, string(respBody))
	}

	fileName := path.Base(targetURL.Path)

	err = ioutil.WriteFile(filepath.Join(destinationDir, fileName), respBody, 0644)
	if err != nil {
		return emptyResponse, err
	}

	err = ioutil.WriteFile(filepath.Join(destinationDir, "version"), []byte(request.Version.ID), 0644)
	if err != nil {
		return emptyResponse, err
	}

	err = ioutil.WriteFile(filepath.Join(destinationDir, "url"), []byte(targetURL.String()), 0644)
	if err != nil {
		return emptyResponse, err
	}

	return Response{
		Version: request.Version,
		Metadata: []initializr.MetadataPair{
			initializr.MetadataPair{
				Name:  "file",
				Value: fileName,
			},
			initializr.MetadataPair{
				Name:  "version",
				Value: request.Version.ID,
			},
		},
	}, nil
}

func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func setValueOrDefault(values *url.Values, key, paramValue, defaultValue string) {
	if empty(paramValue) {
		values.Add(key, defaultValue)
	} else {
		values.Add(key, paramValue)
	}
}

func setValue(values *url.Values, key, paramValue string) {
	if !empty(paramValue) {
		values.Add(key, paramValue)
	}
}
