# Fast Data: Application Logs

TBD

- Estimated time for completion:
 - Install: 20min
 - Development: unbounded
- Target audience: Anyone interested in online log analysis.

**Table of Contents**:

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Install](#install) the demo
- [Use](#use) the demo
- [Development and testing](#development)

## Architecture

![Application Logs demo architecture](img/applogs-architecture.png)

Log data is generated in WordPress (WP) by an end-user interacting with it, this
data gets loaded into Minio and Apache Drill is then used to interactively query it.

## Prerequisites

- A running [DC/OS 1.9.0](https://dcos.io/releases/1.9.0/) or higher cluster with at least 3 private agents and 1 public agent each with 2 CPUs and 5 GB of RAM available as well as the [DC/OS CLI](https://dcos.io/docs/1.8/usage/cli/install/) installed in version 0.14 or higher.
- The [dcos/demo](https://github.com/dcos/demos/) Git repo must be available locally, use: `git clone https://github.com/dcos/demos.git` if you haven't done so, yet.
- The JSON query util [jq](https://github.com/stedolan/jq/wiki/Installation) must be installed.
- [SSH](https://dcos.io/docs/1.9/administration/access-node/sshcluster/) cluster access must be set up.

Going forward we'll call the directory you cloned the `dcos/demo` Git repo into `$DEMO_HOME`.

## Install

### Marathon-LB

For Minio and Apache Drill we need to have Marathon-LB installed:

```bash
$ dcos package install marathon-lb
```

### Minio

To serve the log data for analysis in Drill we use Minio in this demo, just as you would use, say, S3 in AWS.

To set up Minio find out the [IP of the public agent](https://dcos.io/docs/1.9/administration/locate-public-agent/)
and store it in an environment variable called `$PUBLIC_AGENT_IP`, for example:

```bash
$ export PUBLIC_AGENT_IP=52.24.255.200
```

Now you can install the Minio package like so:

```bash
$ cd $DEMO_HOME/1.9/applogs/
$ ./install-minio.sh
```

After this, Minio is available on port 80 of the public agent, so open `$PUBLIC_AGENT_IP`
in your browser and you should see the UI.

Next, we will need to get the Minio credentials in order to access the Web UI (and later on the HTTP API).
The credentials used by Minio are akin to the ones you might know from Amazon S3, called `$ACCESS_KEY_ID`
and `$SECRET_ACCESS_KEY`. In order to obtain these credentials, first go to the `Services` tab of the DC/OS UI and
select the running Minio service. Select the running instance (`minio.535447b6-fcf8-11e6-aa90-8267007bfafa` in my case),
click on the `Files` tab, and select the `stdout` file to view it (by clicking on the magnifying glass icon):

![Minio stdout](img/minio-stdout.png)

Then you should see the output with the credentials:

![Obtaining Minio credentials](img/minio-creds.png)


Note that you can learn more about Minio and the credentials in the respective [example](https://github.com/dcos/examples/tree/master/1.8/minio#using-browser-console).

### Apache Drill

Apache [Drill](http://drill.apache.org/docs/) is a distributed SQL query engine, allowing
you to interactively explore heterogenous datasets across data sources (CSV, JSON, HDFS, HBase, MongoDB, S3).

A prerequisite for the Drill install to work is that three environment variables
are defined: `$PUBLIC_AGENT_IP` (the public agent IP address), as well as `$ACCESS_KEY_ID`
and `$SECRET_ACCESS_KEY` (Minio credentials); all of which are explained in the
previous section. I've been using the following (specific for my setup):

```bash
$ export PUBLIC_AGENT_IP=52.24.255.200
$ export ACCESS_KEY_ID=MRQZLLB72IJRPUGY30MJ
$ export SECRET_ACCESS_KEY=f5nGdq3lxlvpJF1nMOFAgk8h71ZMlM0h4fzUwakj
```

Now do the following to install Drill:

```bash
$ cd $DEMO_HOME/1.9/applogs/
$ ./install-drill.sh
```

Go to `http://$PUBLIC_AGENT_IP:8047/` to access the Drill Web UI:

![Apache Drill Web UI](img/drill-ui.png)

Next we need to configure the S3 storage plugin in order to access data on Minio.
For this, go to the `Storage` tab in Drill, enable the `s3` plugin, click on the `Update` button and paste the content of your (local) [drill-s3-plugin-config.json](drill/drill-s3-plugin-config.json) into the field, overwriting everything which was there in the first place:

![Apache Drill storage plugin config](img/drill-storage-plugin.png)

After another click on the `Update` button the data is stored in ZooKeeper and persists even if you restart Drill.

To check if everything is working fine, go to Minio and create a `test` bucket and upload `drill/apache.log` into it.
Now, go to the Drill UI, change to the `Query` tab and execute the following query to verify your setup:

```sql
select * from s3.`apache.log`
```

You should see something like the following:

![Apache Drill test query result](img/drill-test.png)

### Wordpress

Next we install WordPress, acting as the data source for the logs.
Note that the environment variable called `$PUBLIC_AGENT_IP` must be exported.

```bash
$ cd $DEMO_HOME/1.9/applogs/
$ ./install-wp.sh
```

Discover where WP is available via HAProxy `http://$PUBLIC_AGENT_IP:9090/haproxy?stats`
(look for the `wordpress_XXXXX` frontend):

![WP on Marathon-LB](img/wp-mlb.png)

In my case, WP is available via port `10102` on the public agent, that is via `http://$PUBLIC_AGENT_IP:10102/`:

![WP setup](img/wp.png)

Finally, complete the WP install so that it can be used.

## Use

The following sections describe how to use the demo after having installed it.

First interact with WP, that is, create some posts and surf around. Then, to capture the logs,
execute the following locally (on your machine):

```bash
$ cd $DEMO_HOME/1.9/applogs/
$ echo remote ignore0 ignore1 timestamp request status size origin agent > session.log && dcos task log --lines 1000 wordpress | tail -n +30 | sed 's, \[\(.*\)\] , \"\1\" ,' >> session.log
```

Next upload `session.log` into the `test` bucket in Minio.

Now you can use Drill to understand the usage patterns, for example:

List HTTP requests with more than 1000 bytes payload:

```sql
select remote, request from s3.`session.log` where size > 1000
```

Above query results in something like:

![Apache Drill query result 1](img/query-result1.png)


List HTTP requests that succeeded (HTTP status code 200):

```sql
select remote, request, status from s3.`session.log` where status = 200
```

Above query results in something like:

![Apache Drill query result 2](img/query-result2.png)

## Discussion

TBD.

Should you have any questions or suggestions concerning the demo, please raise an [issue](https://jira.mesosphere.com/) in Jira or let us know via the [users@dcos.io](mailto:users@dcos.io) mailing list.
