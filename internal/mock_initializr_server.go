package internal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

// MockInitializrServer will create a server that mimics a Spring Initializr API
func MockInitializrServer(testdataDir string) *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fileName := request.URL.Path
		if fileName == "/" {
			fileName = "/initializr.json"
		}

		bytes, err := ioutil.ReadFile(filepath.Join(testdataDir, fileName))
		if err != nil {
			response.WriteHeader(500)
			response.Write([]byte(err.Error()))
			return
		}

		response.WriteHeader(200)
		response.Write(bytes)
	}))
}
