#!/bin/bash

docker push mesosphere/dcos-beer-service:latest
docker push mesosphere/dcos-beer-database:latest
docker push mesosphere/dcos-beer-neo4j-migration:latest
docker push mesosphere/dcos-beer-elasticsearch-migration:latest
