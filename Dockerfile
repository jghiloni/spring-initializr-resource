FROM golang:alpine as builder

MAINTAINER Josh Ghiloni <jghiloni@pivotal.io>

COPY . /go/src/github.com/jghiloni/spring-initializr-resource
ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/jghiloni/spring-initializr-resource

RUN dep ensure
RUN go test ./...

RUN go build -o /assets/in ./cmd/in
RUN go build -o /assets/check ./cmd/check

FROM alpine:edge AS resource

# Keep things up to date
RUN apk update
RUN apk upgrade

# Add useful tools for debugging pipelines using this resource type
RUN apk add --no-cache curl bash tzdata ca-certificates unzip zip gzip tar

COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
