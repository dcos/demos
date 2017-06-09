#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

echo deploying Minio ...

: ${PUBLIC_AGENT_IP?"ERROR: The environment variable PUBLIC_AGENT_IP is not set."}

if [ "$(uname)" == "Darwin" ]; then
  # replace the template with the actual value of the public agent IP:
  sed -i '.tmp' "s/PUBLIC_AGENT_IP/$PUBLIC_AGENT_IP/" ./minio-config.json
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  # replace the template with the actual value of the public agent IP:
  sed -i'.tmp' "s/PUBLIC_AGENT_IP/$PUBLIC_AGENT_IP/" ./minio-config.json
fi

# deploy service:
dcos package install minio --options=minio-config.json
# restore template:
mv ./minio-config.json.tmp ./minio-config.json

echo DONE  ====================================================================
