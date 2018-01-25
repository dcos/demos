#!/bin/bash

TAG=1.0

# exit immediately on failure
set -e

# build linux binary
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o beer-service-web-linux

# build docker image
docker build -t mesosphere/dcos-beer-service-web:$TAG .

# push to docker registry
docker push mesosphere/dcos-beer-service-web
