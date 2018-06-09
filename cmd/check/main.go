package main

import (
	"encoding/json"
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
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatalf("reading request from stdin: %s", err.Error())
	}
}

func outputResponse(response check.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("writing response to stdout: %s", err.Error())
	}
}
