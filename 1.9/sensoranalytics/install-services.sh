#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

# pick the first Kafka broker FQDN (note that -r strips the quotes):
broker0=`dcos kafka connection | jq -r .dns[0]`

echo deploying the traffic fetcher ...
# replace the template with the actual value of the broker:
sed -i '.tmp' "s/BROKER/$broker0/" ./service/traffic-fetcher.json
# deploy service:
dcos marathon app add ./service/traffic-fetcher.json
# restore template:
mv ./service/traffic-fetcher.json.tmp ./service/traffic-fetcher.json
echo ==========================================================================

echo deploying the mapping agent ...
# replace the template with the actual value of the broker:
sed -i '.tmp' "s,BROKER,$broker0,; s,_PUBLIC_AGENT_IP,$PUBLIC_AGENT_IP,g; s,_ACCESS_KEY_ID,$ACCESS_KEY_ID,; s,_SECRET_ACCESS_KEY,$SECRET_ACCESS_KEY," ./service/mapping-agent.json
# deploy service:
dcos marathon app add ./service/mapping-agent.json
# restore template:
mv ./service/mapping-agent.json.tmp ./service/mapping-agent.json

echo DONE  ====================================================================
