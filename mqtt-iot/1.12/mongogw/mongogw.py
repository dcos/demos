#!/usr/bin/env python
"""
MQTT to MongoDB Gateway
"""

import json
from argparse import ArgumentParser
import paho.mqtt.client as mqtt
import pymongo
import datetime
import os

parser = ArgumentParser()
parser.add_argument("-b", "--broker", dest="broker_address",
                    required=True, help="MQTT broker address")
parser.add_argument("-p", "--port", dest="broker_port", default=1883, help="MQTT broker port")
parser.add_argument("-m", "--mongouri", dest="mongo_uri", required=True, help="MongoDB URI")
parser.add_argument("-u", "--mongouser", dest="mongo_user", required=True, help="MongoDB user")
parser.add_argument("-w", "--mongopwd", dest="mongo_password", required=True, help="MongoDB password")
args = parser.parse_args()

def on_message(client, userdata, message):
    json_data = json.loads(message.payload)
    post_data = {
        'date': datetime.datetime.utcnow(),
        'deviceUID': json_data['uuid'],
        'value': json_data['value'],
        'gatewayID': os.environ['MESOS_TASK_ID']
    }
    result = devices.insert_one(post_data)
    #print('One device: {0}'.format(result.inserted_id))

# MongoDB connection
mongo_client = pymongo.MongoClient(args.mongo_uri,
                                   username=args.mongo_user,
                                   password=args.mongo_password,
                                   authSource='mongogw',
                                   authMechanism='SCRAM-SHA-1')
db = mongo_client.mongogw
devices = db.devices

# MQTT connection
mqttc = mqtt.Client("mongogw", False)
mqttc.on_message=on_message
mqttc.connect(args.broker_address, args.broker_port)
mqttc.subscribe("device/#", qos=0)
mqttc.loop_forever()
