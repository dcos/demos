#!/usr/bin/env bash

set -o errexit
set -o errtrace 
set -o nounset
set -o pipefail

echo Installing the DC/OS Financial Transaction Processing demo ...

# install InfluxDB
dcos package install --options=./influx-ingest/influx-config.json influxdb --yes

# install Grafana and dependencies
dcos package install marathon-lb --yes
dcos package install grafana --yes

# install Apache Kafka
dcos package install kafka --options=kafka-config.json --yes
echo OK, I will wait 20 sec now to give Kafka some time to get ready ...
sleep 20

# Install the generator and the consumers
bash ./install-services.sh

echo DONE ... now head over to your DC/OS dashboard and configure Grafana with InfluxDB as a datasource