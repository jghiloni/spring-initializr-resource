package in

import "github.com/jghiloni/spring-initializr-resource"

// Request is what will be fed to the in command via stdin
type Request struct {
	Source  initializr.Source  `json:"source"`
	Version initializr.Version `json:"version"`
	Params  Params             `json:"params"`
}

// Params are specified during a 'get' operation and are merged with the Source.
// All are optional
type Params struct {
	Type         string `json:"type,omitempty"`
	Dependencies string `json:"dependencies,omitempty"`
	Packaging    string `json:"packaging,omitempty"`
	JDKVersion   string `json:"jdk_version,omitempty"`
	Language     string `json:"language,omitempty"`
	GroupID      string `json:"group_id,omitempty"`
	ArtifactID   string `json:"artifact_id,omitempty"`
	Version      string `json:"version,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	PackageName  string `json:"package_name,omitempty"`
}

// Response is what is sent back to the container over Stdout
type Response struct {
	Version  initializr.Version        `json:"version"`
	Metadata []initializr.MetadataPair `json:"metadata"`
}
