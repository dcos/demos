#!/usr/bin/env bash

set -o errexit
set -o errtrace 
set -o nounset
set -o pipefail

echo Deploying the fintrans generator
broker0=`dcos kafka connection | jq .dns[0]`
sed -i "s/BROKER/$broker0/g" generator.json
dcos marathon app add generator.json