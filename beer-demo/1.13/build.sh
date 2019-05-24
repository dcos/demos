#!/bin/bash -e

cd $(dirname $0)

cd service
./build.sh
cd ../database
./build.sh
cd ../neo4j-migration
./build.sh
cd ../elasticsearch-migration
./build.sh
