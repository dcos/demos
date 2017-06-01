#!/usr/bin/env bash

set -o errexit
set -o errtrace 
set -o nounset
set -o pipefail

echo Installing the DC/OS Financial Transaction Processing demo ...

# install Apache Kafka
dcos package install kafka --options=kafka-config.json --yes

# install Marathon-LB (needed for Grafana)
dcos package install marathon-lb --yes

# install InfluxDB
dcos package install --options=./influx-ingest/influx-config.json influxdb --yes

# install Grafana and dependencies
dcos package install grafana --yes

echo ==================================================================
echo =
echo = OK, I will wait 2 min now to give Kafka some time to get ready ...
sleep 120

# install the generator and the consumers
echo ==================================================================
echo =
echo = Now the generator and the consumers
bash ./install-services.sh

echo =
echo = DONE
echo ==================================================================

echo Now head over to your DC/OS dashboard and configure Grafana with InfluxDB as a datasource