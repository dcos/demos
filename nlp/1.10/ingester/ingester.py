#!/usr/bin/env python
"""
Load data for demo
"""

import urllib2
import json
import ssl
from kafka import KafkaProducer
import argparse

def main():
    """
    Load JSON from file into Kafka
    """
    parser = argparse.ArgumentParser(description="Load JSON data into Kafka")
    parser.add_argument("-u", "--data_url", dest="data_url", help="URL to fetch JSON data from", required=True)
    parser.add_argument("-b" "--broker", dest="kafka_broker", help="Kafka broker address", required=True)
    parser.add_argument("-t", "--topic", dest="kafka_topic", help="Kafka topic", default="ingest")

    args = parser.parse_args()

    producer = KafkaProducer(value_serializer=lambda m: json.dumps(m).encode('utf-8'),
                             bootstrap_servers=args.kafka_broker)

    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE

    data_file = urllib2.urlopen(data_url, context=ctx)

    for line in data_file:
        producer.send(kafka_topic, json.loads(line))

if __name__ == "__main__":
    main()
