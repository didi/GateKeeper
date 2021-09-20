#!/bin/sh
export GO111MODULE=auto && export GOPROXY=https://goproxy.cn && go mod tidy
GOOS=linux GOARCH=amd64 go build -o ./bin/gatekeeper
docker build -f dockerfile-dashboard -t go-gateteway-dashboard .
docker build -f dockerfile-server -t go-gateteway-server .