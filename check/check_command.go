package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/jghiloni/spring-initializr-resource"
	"github.com/jghiloni/spring-initializr-resource/net"
)

// Command will perform the check operation and look for new versions
type Command struct {
	Client net.HTTPClient
}

// Run will check the specified initializr site and report back new versions from the last check
func (command *Command) Run(request Request) (Response, error) {
	httpResponseBody := struct {
		BootVersion bootVersion `json:"bootVersion"`
	}{}

	req, err := http.NewRequest("GET", request.Source.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.initializr.v2.1+json")

	httpResponse, err := command.Client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("Expected 200 OK, got %d %s with message %s", httpResponse.StatusCode, httpResponse.Status, string(respBody))
	}

	if err = json.Unmarshal(respBody, &httpResponseBody); err != nil {
		return nil, err
	}

	versions := make([]comparableVersion, 0, len(httpResponseBody.BootVersion.Values))
	for _, value := range httpResponseBody.BootVersion.Values {
		buildVersion, releaseType := parseVersion(value)

		if !request.Source.IncludeSnapshots && !strings.EqualFold(releaseType, "RELEASE") {
			continue
		}

		if request.Source.ProductVersion != nil && !request.Source.ProductVersion.Match([]byte(value.ID)) {
			continue
		}

		//value.BuildNum = semver.MustParse(verNum)
		versions = append(versions, comparableVersion{
			version:      value,
			buildVersion: buildVersion,
		})
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].buildVersion.Compare(versions[j].buildVersion) > 0
	})

	if request.Version != nil {
		ver, _ := parseVersion(*request.Version)
		compRequestVersion := comparableVersion{
			version:      *request.Version,
			buildVersion: ver,
		}

		for i := range versions {
			if versions[i].buildVersion.Compare(compRequestVersion.buildVersion) <= 0 {
				return unwrapVersions(versions[:i]), nil
			}
		}
	}

	return unwrapVersions(versions), nil
}

func parseVersion(value initializr.Version) (semver.Version, string) {
	parts := strings.SplitN(value.ID, ".", 4)
	verNum := strings.Join(parts[0:3], ".")
	releaseType := parts[3]

	return semver.MustParse(verNum), releaseType
}

func unwrapVersions(versions []comparableVersion) Response {
	resp := make(Response, 0, len(versions))
	for _, cv := range versions {
		resp = append(resp, cv.version)
	}

	return resp
}

type bootVersion struct {
	Default string               `json:"default"`
	Values  []initializr.Version `json:"values"`
}

type comparableVersion struct {
	version      initializr.Version
	buildVersion semver.Version
}
