#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

: ${PUBLIC_AGENT_IP?"ERROR: The environment variable PUBLIC_AGENT_IP is not set."}
: ${ACCESS_KEY_ID?"ERROR: The environment variable ACCESS_KEY_ID is not set."}
: ${SECRET_ACCESS_KEY?"ERROR: The environment variable SECRET_ACCESS_KEY is not set."}

# pick the first Kafka broker FQDN (note that -r strips the quotes):
broker0=`dcos kafka endpoints broker | jq -r .dns[0]`

echo deploying the traffic fetcher ...
if [ "$(uname)" == "Darwin" ]; then
  # replace the template with the actual value of the broker:
  sed -i '.tmp' "s/BROKER/$broker0/" ./service/traffic-fetcher.json
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  # replace the template with the actual value of the broker:
  sed -i'.tmp' "s/BROKER/$broker0/" ./service/traffic-fetcher.json
fi
# deploy service:
dcos marathon app add ./service/traffic-fetcher.json
# restore template:
mv ./service/traffic-fetcher.json.tmp ./service/traffic-fetcher.json
echo ==========================================================================

echo deploying the mapping agent ...

if [ "$(uname)" == "Darwin" ]; then
  # replace the template with the actual value of the broker:
  sed -i '.tmp' "s,BROKER,$broker0,; s,_PUBLIC_AGENT_IP,$PUBLIC_AGENT_IP,g; s,_ACCESS_KEY_ID,$ACCESS_KEY_ID,; s,_SECRET_ACCESS_KEY,$SECRET_ACCESS_KEY," ./service/mapping-agent.json
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  # replace the template with the actual value of the broker:
  sed -i'.tmp' "s,BROKER,$broker0,; s,_PUBLIC_AGENT_IP,$PUBLIC_AGENT_IP,g; s,_ACCESS_KEY_ID,$ACCESS_KEY_ID,; s,_SECRET_ACCESS_KEY,$SECRET_ACCESS_KEY," ./service/mapping-agent.json
fi
# deploy service:
dcos marathon app add ./service/mapping-agent.json
# restore template:
mv ./service/mapping-agent.json.tmp ./service/mapping-agent.json

echo DONE  ====================================================================
