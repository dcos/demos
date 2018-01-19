#!/bin/bash

cd $(dirname $0)

docker build --tag mesosphere/dcos-beer-database:latest .
