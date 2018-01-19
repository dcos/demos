#!/bin/bash

cd $(dirname $0)

mvn clean install -DskipTests
docker build --tag mesosphere/dcos-beer-service:latest .
