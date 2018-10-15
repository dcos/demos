# Highly Scalable IOT architecture using Mosquitto and MongoDB

This demo presents a highly scalable IOT architecture using Mosquitto, MongoDB and Python microservices on DC/OS. 

- Estimated time for completion: 30min
- Target audience: Anyone interested in deploying IOT architectures.

**Table of Contents**:

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Install](#install)
- [Use the demo](#use)

## Architecture

![Model](images/architecture.jpg)

## Prerequisites

- A running [DC/OS 1.11](https://dcos.io/releases/) or higher cluster with at least 5 private agents and 1 public agent. Each agent should have 2 CPUs and 5 GB of RAM available. The [DC/OS CLI](https://dcos.io/docs/1.11/usage/cli/install/) also needs to be installed.
- The [dcos/demo](https://github.com/dcos/demos/) Git repo must be available locally, use: `git clone https://github.com/dcos/demos/` if you haven't done so yet.
- [SSH](https://docs.mesosphere.com/1.11/administering-clusters/sshcluster/) cluster access must be set up.

The DC/OS services used in the demo are as follows:

- Marathon
- Percona-Server-MongoDB

## Setup

### MongoDB

Our first task is to deploy a MongoDB replica set, we can just use the defaults from the percona-server-mongodb package which will give us a three node set. MongoDB is a widely used choice for deployment in IOT applications due to it's scalability, support for complex queries, and the fact that it doesn’t require strict schemas, so you can push documents to it in native JSON with different types for each field.

Firstly we need to create a config file which sets up the default users, passwords and keys. It should look something like this :

```
Mattbook-Pro:percona matt$ cat demo.json 
{
  "mongodb-credentials": {
    "backupUser": "backup",
    "backupPassword": "backupuserpassword",
    "userAdminUser": "useradmin",
    "userAdminPassword": "useradminpassword",
    "clusterAdminUser": "clusteradmin",
    "clusterAdminPassword": "clusteradminpassword",
    "clusterMonitorUser": "clustermonitor",
    "clusterMonitorPassword": "monitoruserpassword",
    "key": "8cNNTVP6GqEOKzhUVDVryxIt04K6kDbXygamH4upPGAO59gzXVQAgX9NwxwqDvpt 094zMkkRWDLzuCgbg3Aj8EFVEM0/W1Nz+XUSTHEn4HiNzCVG4TTHFP6P1PEPswG6 tQMP6bnRXL7uGWmdGhbAxOV/+p6AfNs67MTvfCeH0EaPCgPPXhJft9D0nZ0SPOm9 VvfxG3djnHClIlclkchoIwc1Kw21loyXwuOjX4RkywVDdmFXjKC+l9yxfyt/9Gyh YE0OlS7ozWLiH8zy0MyzBdK+rc0fsxb2/Kb/8/2diC3O3gdVxjneQxaf66+FHVNW mV9/IHDptBHosdWkv0GboW8ZnTXnk0lyY0Jw85JFuTeFBzqPlB37jR0NU/HFm5QT Ld62woaGIWCTuXGb81QHaglPZUBIhEq/b3tahJBmLc+LKd0FUShoupTtPc2FjxbH xD8dZ+L9Uv7NPtSe+o3sTD60Pnsw1wbOrNDrrC+wpwoMy2GbQjXk/d+SRK/CXfuk Z676GKQDivpinhdF58l4OEi+WEN633yuNtNAQDgz+aOVZKN4oLoyR22B1nrea1qW wzZjRw7kpVxcQKiyn+gDmAZZPbctiVqTNHPE5n9LrOcctuLZKpoQk97lvZTSCKfy d32mfx9szZZ/QCfF9Dt7+G5nJUAULigKnQYRi/i86ZTPHSzfun+ZIzYLCzJuZfyS 7E8DMsmv9wCPrPAF/8cOFMWW0o0Na7GZKCJ8U+AMm92R725h4g5ao6+kQPG7vOkY LR8MJzDOqcmAC0M9AwE5UXQl56V6qBNyREx/WGGYS1B5DOfZvVTJNDkoHVIL1upZ geSlACiXQ+M0Rkgo0h8BJUhGY9LTuc6S8qiMBEnhBClg4kA/u4FJ06nlmF3ZpIXT KsVSr9ee3mu0vSr6P52slvAAX+RL3y+JgSlz2kC8oVgCZZdKn7yq9e6yB3zHNMjX 8VIi/UgFmfqCiaAlUT0pt2ZzGuw1L9QUOuNAZfufSkK1ED4V"
  }
}
```

The passwords for the users must have a minimum of 10 characters, and the key must be at least 1024 characters long, you can generate an appropriate key on MacOS using :

```
Mattbook-Pro:percona matt$ openssl rand -base64 756
```

Now we have our options.json we can use it to install the package :

```
dcos package install percona-server-mongodb --options=demo.json
```

We also need to configure a user for the database we're going to use. To do this we need a JSON file defining the user, which you can find in the demo directory : 

```
Mattbook-Pro:1.11 matt$ cat mongouser.json 
{
  "user": "mongogw",
  "pwd": "123456",
  "roles": [
    { "db": "mongogw", "role": "readWrite" }
  ]
}
```

To add the user :

```
Mattbook-Pro:1.11 matt$ dcos percona-server-mongodb user add mongogw mongouser.json
{
"message": "Received cmd: start update-user with parameters: {MONGODB_CHANGE_USER_DB=mongogw, MONGODB_CHANGE_USER_DATA=eyJ1c2VycyI6W3sidXNlciI6Im1vbmdvZ3ciLCJwd2QiOiIxMjM0NTYiLCJyb2xlcyI6W3sicm9sZSI6InJlYWRXcml0ZSIsImRiIjoibW9uZ29ndyJ9XX1dfQ==}"
}
```

### MQTT layer

Next, we’re going to deploy our MQTT layer. MQTT is a standardised protocol designed for sensors, based around a publish/subscribe model. MQTT producers need to connect to a broker to publish their data, and in our case we’re going to use the open source [Eclipse Mosquitto](https://mosquitto.org/) broker . 

In order to pull from our acquisition layer brokers into our aggregation layer of MongoDB, we’re going to need some kind of gateway, which will subscribe to the relevant topics on our broker, get new messages, reformat them and then write into the MongoDB replica set on the backend. I’ve implemented this as a Python based service, using [paho-mqtt](https://www.eclipse.org/paho/clients/python/) and the [pymongo](https://api.mongodb.com/python/current/) library.

Since we may theoretically have many thousands of devices, it’s important our acquisition layer is scalable. We’re going to achieve this on the MQTT side by load balancing multiple instances of Mosquitto endpoints using a named VIP in DC/OS, although in a real world internet facing example we would probably use Marathon-LB instances exposed outside of our cluster. 

Whilst this will work for the device facing side, if we were to use the VIP from our gateway layer, and we want that also to be scalable, then we could potentially end up with missing data since our gateways could be reading from any of the Mosquitto instances, and the MQTT layer isn’t clustered. To solve this problem, we’ll deploy each Mosquitto instance together with it’s own dedicated gateway microservice in a pod, and have the gateway microservice connect to Mosquitto over the localhost connection, ensuring each gateway only gets data from a single Mosquitto instance. So for a device connecting via the loadbalancer address, it can connect to any of our deployed Mosquitto instances and our storage layer will get the data.

The gateway microservice is in the mongogw folder in the demo repository, and we will deploy it in a Docker container. The Dockerfile which was used to build the container is also available in the repository. Let’s look quickly at the [code](https://github.com/dcos/demos/blob/master/mqtt-iot/1.11/mongogw/mongogw.py) for this microservice.

In this case, we connect to both the MQTT broker and MongoDB, then we subscribe to all of the MQTT topics prefixed with device/ and then use a callback function each time we read a message, to reformat the message and add some additional metadata of our own, before writing it into MongoDB. We add a timestamp at this stage, which reflects the possibility that in a real world example our sensor devices may not have enough processing power to provide a clock function. We also add the Mesos task ID of the gateway process, which will allow us to track which gateway the message came through as we scale up the acquisition layer. 

To deploy the pod, we use the mqttpod.json in the root folder of the demo repository :

```
dcos marathon pod add mqttpod.json
```

### Devices

Now our acquisition layer is up and running, we can start our devices. For the purposes of this demo, the devices are just emulated as a python microservice, with a single value generator publishing via MQTT to the acquisition layer. As with the gateway microservice, we will deploy this in a Docker container. The code for the device and the Dockerfile used to build the container is in the device folder in the demo repository. The [code](https://github.com/dcos/demos/blob/master/mqtt-iot/1.11/device/device.py) for this is very simple.

As we can see, our device takes a few different arguments. Firstly the address and port of the MQTT broker to connect to, the rate at which we want to generate values, and a currently unused QOS setting. QOS ( Quality of Service ) is part of the MQTT standard, and you can learn more about this in the Mosquitto [docs](https://mosquitto.org/man/mqtt-7.html).

On startup the device generates a unique UUID, connects to the MQTT broker on a topic of device/$uuid, and then every sample cycle it generates a random number within a range and then publishes a small JSON payload containing it’s UUID and the value. 


To deploy this, from the root of the demo repository :

```
dcos marathon app add device/device.json
```

Once our first device is up and running, we can check in MongoDB to ensure all of our layers are operating correctly. First we need a task ID for one of our MongoDB replica set :

```
Mattbook-Pro:1.11 matt$ dcos task
NAME                 HOST        USER  STATE  ID                                                            MESOS ID                                 REGION  ZONE  
device               10.0.2.202  root    R    device.612f3f82-850f-11e8-9832-a63d4d1f1cb3                   207783f7-10cf-44ab-9966-9bb1d6f3b797-S8   ---    ---   
mongo-rs-0-mongod    10.0.0.26   root    R    mongo-rs-0-mongod__1b3f7f6c-ea11-4848-a57b-e22b68cb46e7       207783f7-10cf-44ab-9966-9bb1d6f3b797-S0   ---    ---   
mongo-rs-1-mongod    10.0.0.241  root    R    mongo-rs-1-mongod__a68630cb-a6d8-4c61-b83b-45108363df1f       207783f7-10cf-44ab-9966-9bb1d6f3b797-S3   ---    ---   
mongo-rs-2-mongod    10.0.2.168  root    R    mongo-rs-2-mongod__55e58da7-cd56-461e-a9c5-ff8a1bb07ed2       207783f7-10cf-44ab-9966-9bb1d6f3b797-S5   ---    ---   
mongogw              10.0.0.26   root    R    mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw    207783f7-10cf-44ab-9966-9bb1d6f3b797-S0   ---    ---   
mosquitto            10.0.0.26   root    R    mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mosquitto  207783f7-10cf-44ab-9966-9bb1d6f3b797-S0   ---    ---   
percona-mongo        10.0.0.241  root    R    percona-mongo.b468ffe6-8509-11e8-9832-a63d4d1f1cb3            207783f7-10cf-44ab-9966-9bb1d6f3b797-S3   ---    ---   
watchdog-0-watchdog  10.0.2.202  root    R    watchdog-0-watchdog__f33ccb02-ef87-4bca-86db-98c4382d19ef     207783f7-10cf-44ab-9966-9bb1d6f3b797-S8   ---    ---   
```

Once we have a task ID, we can use the DC/OS CLI to get a shell inside that container :

```
Mattbook-Pro:1.11 matt$ dcos task exec --tty --interactive mongo-rs-0-mongod__1b3f7f6c-ea11-4848-a57b-e22b68cb46e7 /bin/bash
root@ip-10-0-0-26:/mnt/mesos/sandbox# 
```

And then use the mongo shell to connect to MongoDB :

```
root@ip-10-0-0-26:/mnt/mesos/sandbox# mongo mongodb://mongogw:123456@mongo-rs-0-mongod.percona-mongo.autoip.dcos.thisdcos.directory,mongo-rs-1-mongod.percona-mongo.autoip.dcos.thisdcos.directory,mongo-rs-2-mongod.percona-mongo.autoip.dcos.thisdcos.directory:27017/mongogw?replicaSet=rs
```

Note that we’re using the user we created earlier, along with the DC/OS automatically assigned DNS entries for each of our MongoDB replica set. 

Now we should end up with a MongoDB shell prompt, so let’s switch to our database :

```
rs:PRIMARY> use mongogw;
switched to db mongogw
```

From here, we can count how many entries we now have in our devices collection :

```
rs:PRIMARY> db.devices.count();
117
```

And we can also look at the inserted data :

```
rs:PRIMARY> db.devices.findOne();
{
	"_id" : ObjectId("5b46066a8c454b00112a0bb9"),
	"date" : ISODate("2018-07-11T13:30:18.066Z"),
	"deviceUID" : "21366038-3d21-415f-a265-fa3499e15559",
	"value" : 24,
	"gatewayID" : "mqtt.instance-782e4fbd-850d-11e8-9832-a63d4d1f1cb3.mongogw"
}
```

From this output, we can see the re-formatted data coming in from our gateway process, which has also added a timestamp, and the task ID of the gateway process in DC/OS.

So we can now see data being acquired for our single device, let’s scale up the devices. In our device.json, let’s increase the number to 3 :

```
(venv) Mattbook-Pro:device matt$ cat device.json 
{
  "id": "device",
  "instances": 3,
  "cpus": 0.1,
  "mem": 16,
  "cmd": "./device.py -b mqtt.marathon.l4lb.thisdcos.directory -r 2",
  "container": {
    "type": "MESOS",
    "docker": {
      "image": "mattjarvis/device",
      "forcePullImage": true,
      "privileged": false
    }
  },
  "requirePorts": false
}
```

And now we’ll use that modified JSON to update our app :

```
dcos marathon app update device < device/device.json
```

We can now check in the DC/OS Services tab and see the additional instances of the device app being deployed :

![Device Scaling](images/devicescaled.png)

We can also now see different UUID’s in the data being received by our MongoDB layer :

```
rs:PRIMARY> db.devices.find().limit(5).sort({$natural:-1})
{ "_id" : ObjectId("5b461a613613f800047dcaf2"), "date" : ISODate("2018-07-11T14:55:29.323Z"), "deviceUID" : "15ea39c3-26ba-448d-809c-19112412b4a5", "value" : 21, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461a613613f800047dcaf1"), "date" : ISODate("2018-07-11T14:55:29.183Z"), "deviceUID" : "72aef257-9621-4348-bd24-d1eda173f4f7", "value" : 25, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461a613613f800047dcaf0"), "date" : ISODate("2018-07-11T14:55:29.066Z"), "deviceUID" : "cd3da0de-14b3-423f-a1bb-ef64173f29bb", "value" : 29, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461a5f3613f800047dcaef"), "date" : ISODate("2018-07-11T14:55:27.321Z"), "deviceUID" : "15ea39c3-26ba-448d-809c-19112412b4a5", "value" : 25, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461a5f3613f800047dcaee"), "date" : ISODate("2018-07-11T14:55:27.180Z"), "deviceUID" : "72aef257-9621-4348-bd24-d1eda173f4f7", "value" : 27, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
```

As we can see, these are all using the same gateway, as we only have one deployed, so let’s go and deploy another pod in that layer. Edit mqttpod.json to use :

```
 "instances" : 2
```

and now update the Marathon configuration for the running pod :

```
dcos marathon pod update mqtt < mqttpod.json
```

Since all of our devices are already attached to the first gateway through the load balanced VIP, let’s increase the device count to 5 instances. Edit the device.json again to use :

```
 "instances": 5
```

and now update the Marathon configuration for the device application :

```
dcos marathon app update device < device/device.json
```

Once our new devices are up and running, because they are using the load balancer, we should see some of them connect to the new gateway. We can query MongoDB again and see we now have more than one gateway ID in our received data :

```
rs:PRIMARY> db.devices.find().limit(5).sort({$natural:-1})
{ "_id" : ObjectId("5b461b834c3f740004d684ef"), "date" : ISODate("2018-07-11T15:00:19.893Z"), "deviceUID" : "8a69d293-22fe-4e00-9641-b3c04ad6b011", "value" : 30, "gatewayID" : "mqtt.instance-b3e911f5-851a-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461b834c3f740004d684ee"), "date" : ISODate("2018-07-11T15:00:19.687Z"), "deviceUID" : "d1e4a552-e926-4897-8fab-0148e504ca49", "value" : 21, "gatewayID" : "mqtt.instance-b3e911f5-851a-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461b833613f800047dcca5"), "date" : ISODate("2018-07-11T15:00:19.642Z"), "deviceUID" : "15ea39c3-26ba-448d-809c-19112412b4a5", "value" : 20, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461b833613f800047dcca4"), "date" : ISODate("2018-07-11T15:00:19.491Z"), "deviceUID" : "72aef257-9621-4348-bd24-d1eda173f4f7", "value" : 20, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
{ "_id" : ObjectId("5b461b833613f800047dcca3"), "date" : ISODate("2018-07-11T15:00:19.396Z"), "deviceUID" : "cd3da0de-14b3-423f-a1bb-ef64173f29bb", "value" : 21, "gatewayID" : "mqtt.instance-535430a1-850f-11e8-9832-a63d4d1f1cb3.mongogw" }
```

Finally we can also scale our MongoDB layer, horizontally by adding further instances to our replica set or vertically by changing the size of the instances. For the purposes of this demo, let’s add another two instances to our replica set. We can do this very simply using the DC/OS CLI :

```
Mattbook-Pro:percona matt$ dcos percona-server-mongodb scale up 5
Pod task count for 'percona-server-mongodb' is currently 3
Updating service configuration for 'percona-server-mongodb'
Pod task count for 'percona-server-mongodb' is updated to 5
```

Once we’ve changed the configuration, click on Review and Run, then Run Service, and the percona-mongo service will deploy an additional two instances to our replica set. 

Once the deployment is completed, we’ll see another two instances of mongod in the Services tab, with no interruptions to the MongoDB service itself. 

![Scale MongoDB](images/mongoscaled.png)

So we can see that we’ve ended up with a highly scalable architecture in each of the layers of our demo IOT application, and DC/OS made it extremely easy to deploy and manage. As the number of devices increases, we can simply deploy more Mosquitto/Gateway pods behind our load balancer, and if we need additional capacity or performance on our storage layer, we can very easily scale the MongoDB layer using the percona-mongo framework to either scale vertically or horizontally.

### 

Should you have any questions or suggestions concerning the demo, please raise an [issue](https://jira.mesosphere.com/) in Jira or let us know via the [users@dcos.io](mailto:users@dcos.io) mailing list.

