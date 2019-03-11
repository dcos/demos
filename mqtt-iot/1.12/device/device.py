#!/usr/bin/env python
"""
MQTT generator
"""

import random
import time
import uuid
import json
from argparse import ArgumentParser
import paho.mqtt.client as mqtt

parser = ArgumentParser()
parser.add_argument("-b", "--broker", dest="broker_address",
                    required=True, help="MQTT broker address")
parser.add_argument("-p", "--port", dest="broker_port", default=1883, help="MQTT broker port")
parser.add_argument("-r", "--rate", dest="sample_rate", default=5, help="Sample rate")
parser.add_argument("-q", "--qos", dest="qos", default=0, help="MQTT QOS")
args = parser.parse_args()

uuid = str(uuid.uuid4())
topic = "device/%s" % uuid

mqttc = mqtt.Client(uuid, False)
mqttc.connect(args.broker_address, args.broker_port)
while True:
    rand = random.randint(20,30)
    msg = {
            'uuid': uuid,
            'value': rand
    }
    mqttc.publish(topic, payload=json.dumps(msg), qos=args.qos)
    time.sleep(float(args.sample_rate))
mqttc.loop_forever()
