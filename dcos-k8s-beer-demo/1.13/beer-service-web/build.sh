#!/usr/bin/env bash
set -e

TAG=1.13

# build linux binary
make

# build docker image
docker build -t dcoslabs/dcos-beer-service-web:$TAG .

# push to docker registry
docker push dcoslabs/dcos-beer-service-web:$TAG
