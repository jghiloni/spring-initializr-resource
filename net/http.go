package net

import "net/http"

// HTTPClient is an interface that can be used in place of http.Client and is suitable for testing
//go:generate counterfeiter . HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
