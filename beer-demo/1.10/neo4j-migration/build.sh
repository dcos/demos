#!/bin/bash

cd $(dirname $0)

mvn clean install
docker build --tag mesosphere/dcos-beer-neo4j-migration:latest .
