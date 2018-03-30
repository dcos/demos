#!/usr/bin/env bash
HOSTNAME=$1
HOSTPORT=$2
if [ -z "$HOSTNAME" ];
then
    echo "**************************************************";
    echo "please specify hostname + port (optional):";
    echo "";
    echo "cassandra-dcos-node.cassandra.dcos.mesos 9042";
    echo "";
    echo "**************************************************";
    exit 1
fi
cd /opt/bus-demo/commons/src/main/resources
cqlsh $HOSTNAME $HOSTPORT < create_tables.cql
