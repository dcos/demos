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

## Prerequisites

- A running [DC/OS 1.8.7](https://dcos.io/releases/1.8.7/) or higher cluster with at least 4 private agents and 1 public agent each with 2 CPUs and 5 GB of RAM available as well as the [DC/OS CLI](https://dcos.io/docs/1.8/usage/cli/install/) installed in version 0.14 or higher.
- The [dcos/demo](https://github.com/dcos/demos/) Git repo must be available locally, use: `git clone https://github.com/dcos/demos.git` if you haven't done so, yet.
- The JSON query util [jq](https://github.com/stedolan/jq/wiki/Installation) must be installed.
- [SSH](https://dcos.io/docs/1.8/administration/access-node/sshcluster/) cluster access must be set up.

Going forward we'll call the directory you cloned the `dcos/demo` Git repo into `$DEMO_HOME`.

## Install

### Wordpress

TBD

### Minio

To serve some static data we use Minio in this demo, just as you would use, say, S3 in AWS.

In order to use Minio you first need to have Marathon-LB installed:

```bash
$ dcos package install marathon-lb
```

Next find out the [IP of the public agent](https://dcos.io/docs/1.8/administration/locate-public-agent/)
and store it in an environment variable called `$PUBLIC_AGENT_IP`, for example:

```bash
$ export PUBLIC_AGENT_IP=34.250.247.12
```

Now you can install the Minio package like so:

```bash
$ cd $DEMO_HOME/1.8/sensoranalytics/
$ sed -i '.tmp' "s/PUBLIC_AGENT_IP/$PUBLIC_AGENT_IP/" ./minio-config.json
$ dcos package install minio --options=minio-config.json
$ mv ./minio-config.json.tmp ./minio-config.json
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
$ sed -i '.tmp' "s,_PUBLIC_AGENT_IP,$PUBLIC_AGENT_IP,g; s,_ACCESS_KEY_ID,$ACCESS_KEY_ID,; s,_SECRET_ACCESS_KEY,$SECRET_ACCESS_KEY," ./service/apache-drill.json
$ dcos package marathon app add ./service/apache-drill.json
$ mv ./service/apache-drill.json.tmp ./service/apache-drill.json
```

Got to `http://$PUBLIC_AGENT_IP:8047/` to access the Drill Web UI:

![Apache Drill Web UI](img/drill-ui.png)

## Use

The following sections describe how to use the demo after having installed it.

## Development

If you are interested in testing this demo locally or want to extend it, follow the instructions in this section.

### Tunneling

For local development and testing we use [DC/OS tunneling](https://dcos.io/docs/1.8/administration/access-node/tunnel/) to make the nodes directly accessible on the development machine. The following instructions are only an example (using Tunnelblick on macOS) and the concrete steps necessary depend on your platform as well as on what VPN client you're using.

```bash
$ sudo dcos tunnel vpn --client=/Applications/Tunnelblick.app/Contents/Resources/openvpn/openvpn-2.3.12/openvpn
```

Note that it is necessary to [add the announced DNS servers]( https://support.apple.com/kb/PH18499?locale=en_US) as told by Tunnelblick, and make sure the are they appear at the top of the list, before any other DNS server entries.

### Apache Drill

Use `drill-s3-plugin-config.json` to configure the S3 storage plugin.

Execute:

```sql
select * from s3.`apache.log`
```

## Discussion



Should you have any questions or suggestions concerning the demo, please raise an [issue](https://dcosjira.atlassian.net/) in Jira or let us know via the [users@dcos.io](mailto:users@dcos.io) mailing list.
