#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

# Install the required topics
echo "Installing topics"
input_topics=(London SF Tokyo Moscow NYC)
output_topics=(dollarrate maxdollars ratecount)
for topic in "${input_topics[@]}"
do
    dcos kafka topic create $topic --partitions 1 --replication 1 || exit 1
done
for topic in "${output_topics[@]}"
do
    dcos kafka topic create $topic --partitions 1 --replication 1 || exit 1
done
# pick the first Kafka broker FQDN (note that -r strips the quotes):
broker0=`dcos kafka connection | jq -r .dns[0]`

echo deploying the fintrans generator ...
if [ "$(uname)" == "Darwin" ]; then
  # replace the template with the actual value of the broker:
  sed -i '.tmp' "s/BROKER/$broker0/" generator.json
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  # replace the template with the actual value of the broker:
  sed -i'.tmp' "s/BROKER/$broker0/" generator.json
fi
# deploy service:
dcos marathon app add generator.json
# restore template:
mv generator.json.tmp generator.json
echo ==========================================================================

