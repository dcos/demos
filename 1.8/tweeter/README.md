# Tweeter

Tweeter is a sample service that demonstrates how easy it is to run a Twitter-like service on DC/OS.

Capabilities:

* Stores tweets in Cassandra
* Streams tweets to Kafka as they come in
* Real time tweet analytics with Spark and Zeppelin

## Install Prerequisites on Your Machine

To run the demo script (`cli_script.sh`), the following pieces of software are expected to be available:

* dcos cli
* cypress

### Installing cypress

Cypress is a nodejs package that executes UI tests. This is used by the demo script to submit a tweet and verify it is displayed within the Tweeter UI.

To install on OSX, perform the following from within this project's directory:

```
$ brew update
$ brew install node
$ npm install -g cypress-cli
$ cypress install
```

You can run the UI tests separately with `cypress run`.

### Configuring cypress

The demo script has a section that creates the JSON file `ci-conf.json`. This file is read by the cypress to determine the URL of the DC/OS cluster, the URL of the tweeter application, and the log in credentials to use. Without this file the UI tests will fail.

Example:

```
{
  "tweeter_url": "52.xx.xx.xx:10000",
  "url": "http://my-cool-demo.us-west-2.elb.amazonaws.com/",
  "username": "admin",
  "password": "password"
}
```

## Demo Cluster Prerequisites

You'll need a DC/OS cluster with one public node and at least five private nodes and the DC/OS CLI locally installed.

## Scripted Demo

`cli_script.sh` in this repository can be utilized to setup a tweeter demo
cluster automatically and provide some random traffic, however the GUI experience with Zeppelin analytics and DC/OS feature presentation must still be done by hand:
* Install Zeppelin from the GUI using the default values
* Log into the Tweeter UI at http://[elb_hostname] and post a sample tweet
* Start the tweeter load job from the CLI using the command dcos/bin/dcos marathon app add post-tweets.json
* Kill one of the Tweeter containers in Marathon and show that the Tweeter is still up and tweets are still flowing in
* Log into Zeppelin using the https interface at https://[master_ip]/service/zeppelin
* Click Import note and import tweeter-analytics.json from the Tweeter repo clone you made locally
* Open the newly loaded "Tweeter Analytics" Notebook
* Run the Load Dependencies step to load the required libraries into Zeppelin
* Run the Spark Streaming step, which reads the tweet stream from Zookeeper, and puts them into a temporary table that can be queried using SparkSQL - this spins up the Zeppelin spark context so you can show them the increased utilization on the dashboard
* Next, run the Top Tweeters SQL query, which counts the number of tweets per user, using the table created in the previous step
* The table updates continuously as new tweets come in, so re-running the query will produce a different result every time


### Stage EBC Demo
Run Tweeter against an EE cluster, but do not start Zeppelin or post tweets
```
$ bash cli_script.sh --infra --url http://my.dcos.url
```

### Get Manual Demo script and Run Nothing
This combination of options will not actually run the demo but stop and prompt you with the appropriate CLI command to do to run the demo. Note, by specifying -oss without a DCOS\_AUTH\_TOKEN set, the dummy CI auth token will be used
```
$ bash cli_script.sh --step --manual --oss --url http://my.dcos.url
```

### Open DC/OS Tweeter Demo Setup
The steps below are applicable for Open DC/OS, when it does not have a super
set. The auth token we set below

#### Login to dcos and copy the auth token:
```
$ dcos auth login

Please go to the following link in your browser:

    http://54.70.182.15/login?redirect_uri=urn:ietf:wg:oauth:2.0:oob

Enter authentication token:

```

If you wish to access this token again later, use the cli command:
```
dcos config show core.dcos_acs_token
```

#### Set DCOS Auth Token to the environment variable

```
export DCOS_AUTH_TOKEN=<TOKEN_FROM_PREVIOUS_STEP>
```
#### Run the cli script
```
$ ./cli_script.sh --oss --url http://52.70.182.15
```

## Manual Test Steps
### Install the cluster prereqs
```
$ dcos package install marathon-lb
$ dcos package install cassandra --cli
$ dcos package install kafka --cli
$ dcos package install zeppelin
```

Wait until the Kafka and Cassandra services are healthly. You can check their status with:

```
$ dcos kafka connection
...
$ dcos cassandra connection
...
```

### Edit the Tweeter Service Config

Edit the `HAPROXY_0_VHOST` label in `tweeter.json` to match your public ELB hostname. Be sure to remove the leading `http://` and the trailing `/` For example:

```json
{
  "labels": {
    "HAPROXY_0_VHOST": "brenden-7-publicsl-1dnroe89snjkq-221614774.us-west-2.elb.amazonaws.com"
  }
}
```

### Run the Tweeter Service

Launch three instances of Tweeter on Marathon using the config file in this repo:

```
$ dcos marathon app add tweeter.json
```

The service talks to Cassandra via `node-0.cassandra.mesos:9042`, and Kafka via `broker-0.kafka.mesos:9557` in this example.

Traffic is routed to the service via marathon-lb. Navigate to `http://<public_elb>` to see the Tweeter UI and post a Tweet.


### Post a lot of Tweets

Post a lot of Shakespeare tweets from a file:

```
dcos marathon app add post-tweets.json
```

This will post more than 100k tweets one by one, so you'll see them coming in steadily when you refresh the page. Take a look at the Networking page on the UI to see the load balancing in action.


### Streaming Analytics

Next, we'll do real-time analytics on the stream of tweets coming in from Kafka.

Navigate to Zeppelin at `https://<master_public_ip>/service/zeppelin/`, click `Import note` and import `tweeter-analytics.json`. Zeppelin is preconfigured to execute Spark jobs on the DCOS cluster, so there is no further configuration or setup required.

Run the *Load Dependencies* step to load the required libraries into Zeppelin. Next, run the *Spark Streaming* step, which reads the tweet stream from Zookeeper, and puts them into a temporary table that can be queried using SparkSQL. Next, run the *Top Tweeters* SQL query, which counts the number of tweets per user, using the table created in the previous step. The table updates continuously as new tweets come in, so re-running the query will produce a different result every time.


NOTE: if /service/zeppelin is showing as Disconnected (and hence can’t load the notebook), make sure you're using HTTPS instead of HTTP, until [this PR](https://github.com/dcos/dcos/pull/27) gets merged. Alternatively, you can use marathon-lb. To do this, add the following labels to the Zeppelin service and restart:


`HAPROXY_0_VHOST = [elb hostname]`

`HAPROXY_GROUP = external`

You can get the ELB hostname from the CCM “Public Server” link.  Once Zeppelin restarts, this should allow you to use that link to reach the Zeppelin GUI in “connected” mode.

## Developing Tweeter

You'll need Ruby and a couple of libraries on your local machine to hack on this service. If you just want to run the demo, you don't need this.

### Homebrew on Mac OS X

Using Homebrew, install `rbenv`, a Ruby version manager:

```bash
$ brew update
$ brew install rbenv
```

Run this command and follow the instructions to setup your environment:

```bash
$ rbenv init
```

To install the required Ruby version for Tweeter, run from inside this repo:

```bash
$ rbenv install
```

Then install the Ruby package manager and Tweeter's dependencies. From this repo run:

```bash
$ gem install bundler
$ bundle install
```
