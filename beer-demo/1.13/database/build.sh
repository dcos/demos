#!/usr/bin/env bash

cd $(dirname $0)

docker build --tag dcoslabs/dcos-beer-database:1.13 .
