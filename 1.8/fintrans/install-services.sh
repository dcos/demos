#!/usr/bin/env bash

set -o errexit
set -o errtrace 
set -o nounset
set -o pipefail

# pick the first Kafka broker FQDN (note that -r strips the quotes):
broker0=`dcos kafka connection | jq -r .dns[0]`

echo deploying the fintrans generator ...
# replace the template with the actual value of the broker:
sed -i '.tmp' "s/BROKER/$broker0/" ./service/generator.json
# deploy service:
dcos marathon app add ./service/generator.json
# restore template:
mv ./service/generator.json.tmp ./service/generator.json
echo ==========================================================================

echo deploying the recent financial transactions consumer ...
# replace the template with the actual value of the broker:
sed -i '.tmp' "s/BROKER/$broker0/" ./service/influx-ingest.json
# deploy service:
dcos marathon app add ./service/influx-ingest.json
# restore template:
mv ./service/influx-ingest.json.tmp ./service/influx-ingest.json
echo ==========================================================================

echo deploying the money laundering detector ...
# replace the template with the actual value of the broker:
sed -i '.tmp' "s/BROKER/$broker0/" ./service/laundering-detector.json
# deploy service:
dcos marathon app add ./service/laundering-detector.json
# restore template:
mv ./service/laundering-detector.json.tmp ./service/laundering-detector.json

echo DONE  ====================================================================
