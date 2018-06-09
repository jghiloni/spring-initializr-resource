package check

import (
	"github.com/jghiloni/spring-initializr-resource"
)

// Request is the input that will be fed to the check command
type Request struct {
	Source  initializr.Source   `json:"source"`
	Version *initializr.Version `json:"version,omitempty"`
}

// Response is the output of the check command
type Response []initializr.Version
