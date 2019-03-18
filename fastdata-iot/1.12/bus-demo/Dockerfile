# Based on https://github.com/ANierbeck/BusFloatingData/tree/master/terraform/smack/data
FROM cassandra:3.11.4

RUN apt-get update && \
    apt-get install -y git && \
    mkdir -p /opt/bus-demo

ADD import_data.sh /opt/bus-demo/import_data.sh
ADD create_tables.cql /opt/bus-demo/create_tables.cql

ENTRYPOINT ["/opt/bus-demo/import_data.sh"]
