#!/bin/bash

# exit immediately on failure
set -e

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o fraudDisplay-linux

