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

- A running [DC/OS 1.8.7](https://dcos.io/releases/1.8.7/) or higher cluster with at least 3 private agents and 1 public agent each with 2 CPUs and 5 GB of RAM available as well as the [DC/OS CLI](https://dcos.io/docs/1.8/usage/cli/install/) installed in version 0.14 or higher.
- The [dcos/demo](https://github.com/dcos/demos/) Git repo must be available locally, use: `git clone https://github.com/dcos/demos.git` if you haven't done so, yet.
- The JSON query util [jq](https://github.com/stedolan/jq/wiki/Installation) must be installed.
- [SSH](https://dcos.io/docs/1.8/administration/access-node/sshcluster/) cluster access must be set up.

Going forward we'll call the directory you cloned the `dcos/demo` Git repo into `$DEMO_HOME`.

## Install

### Marathon-LB

For Minio and Apache Drill we need to have Marathon-LB installed:

```bash
$ dcos package install marathon-lb
```

### Minio

To serve the log data for analysis in Drill we use Minio in this demo, just as you would use, say, S3 in AWS.

To set up Minio find out the [IP of the public agent](https://dcos.io/docs/1.8/administration/locate-public-agent/)
and store it in an environment variable called `$PUBLIC_AGENT_IP`, for example:

```bash
$ export PUBLIC_AGENT_IP=34.250.247.12
```

Now you can install the Minio package like so:

```bash
$ cd $DEMO_HOME/1.8/applogs/
$ ./install-minio.sh
```

After this, Minio is available on port 80 of the public agent, so open `$PUBLIC_AGENT_IP`
in your browser and you should see the UI.

Next, we will need to get the Minio credentials in order to access the Web UI (and later on the HTTP API).
The credentials used by Minio are akin to the ones you might know from Amazon S3, called `$ACCESS_KEY_ID`
and `$SECRET_ACCESS_KEY`. In order to obtain these credentials, go to the `Services` tab of the DC/OS UI and
select the running Minio service; click on the `Logs` tab and you should see:

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
$ export PUBLIC_AGENT_IP=34.250.247.12
$ export ACCESS_KEY_ID=F3QE89J9WPSC49CMKCCG
$ export SECRET_ACCESS_KEY=2/parG/rllluCLMgHeJggJfY9Pje4Go8VqOWEqI9
```

Now do the following to install Drill:

```bash
$ cd $DEMO_HOME/1.8/applogs/
$ ./install-drill.sh
```

Go to `http://$PUBLIC_AGENT_IP:8047/` to access the Drill Web UI:

![Apache Drill Web UI](img/drill-ui.png)

Next we need to configure the S3 storage plugin in order to access data on Minio.
For this, go to the `Storage` tab in Drill, enable the `s3` plugin, click on the `Update` button and paste the content of your [drill-s3-plugin-config.json](drill/drill-s3-plugin-config.json) into the field, overwriting everything which was there in the first place:

![Apache Drill storage plugin config](img/drill-storage-plugin.png)

After another click on the `Update` button the data is stored in ZooKeeper and persists even if you restart Drill.

To check if everything is working fine, create a `test` bucket and upload `drill/apache.log` into it
Execute the following query to verify your setup:

```sql
select * from s3.`apache.log`
```

### Wordpress

Note that the environment variable called `$PUBLIC_AGENT_IP` must be exported.

```bash
$ cd $DEMO_HOME/1.8/applogs/
$ ./install-wp.sh
```

Discover where WP is available via HAProxy `http://$PUBLIC_AGENT_IP:9090/haproxy?stats`:

![WP on Marathon-LB](img/wp-mlb.png)

Finally, complete the WP install so that it can be used.

## Use

The following sections describe how to use the demo after having installed it.

First interact with WP, create posts and surf around. Then, to capture the logs,
execute the following locally (on your machine):

```bash
$ echo remote ignore0 ignore1 timestamp request status size origin agent > session.log && dcos task log --lines 1000 wordpress | tail -n +5 | sed 's, \[\(.*\)\] , \"\1\" ,' >> session.log
```

Next upload `session.log` into the `test` bucket in Minio.

Now you can use Drill to understand the usage patterns, for example:

```sql
select remote, request from s3.`session.log` where size > 1000

select remote, request, status from s3.`session.log` where size > 1000 AND status = 200
```

Result:

![Apache Drill query result](img/query-result.png)

## Discussion

TBD.

Should you have any questions or suggestions concerning the demo, please raise an [issue](https://dcosjira.atlassian.net/) in Jira or let us know via the [users@dcos.io](mailto:users@dcos.io) mailing list.
