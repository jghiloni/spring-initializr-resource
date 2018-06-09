package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/jghiloni/spring-initializr-resource"
	"github.com/jghiloni/spring-initializr-resource/check"
)

func main() {
	var request check.Request
	inputRequest(&request)

	client, err := initializr.NewHTTPClient(request.Source)
	if err != nil {
		log.Fatalf("error creating HTTP client: %s", err.Error())
	}

	command := &check.Command{
		Client: client,
	}

	response, err := command.Run(request)
	if err != nil {
		log.Fatal(err)
	}

	outputResponse(response)
}

func inputRequest(request *check.Request) {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("reading request from stdin: %s", err.Error())
	}

	if f, err := ioutil.TempFile(os.TempDir(), "check-request-"); err != nil {
		log.Printf("could not log request from stdin but will continue anyway: %s", err.Error())
	} else {
		defer f.Close()
		f.Write(stdin)
	}

	if err := json.Unmarshal(stdin, request); err != nil {
		log.Fatalf("decoding request: %s", err.Error())
	}
}

func outputResponse(response check.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("writing response to stdout: %s", err.Error())
	}
}
