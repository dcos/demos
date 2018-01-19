#!/bin/bash

cd $(dirname $0)

mvn clean install
docker build --tag mesosphere/dcos-beer-elasticsearch-migration:1.0 .
