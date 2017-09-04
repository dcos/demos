# Tweeter

In this tutorial you install and deploy a containerized Ruby on Rails app named Tweeter. Tweeter is an app similar to Twitter that you can use to post 140-character messages to the internet. Then, you use Zeppelin to perform real-time analytics on the data created by Tweeter.

- Estimated time for completion:
 - Manual: 25min
- Target audience: Anyone interested microservices and stream data processing with Apache Spark, Apache Cassandra, and Apache Kafka.

In this tutorial you will deploy these (micro)services to your DC/OS cluster.

- The DC/OS [Cassandra][1] service is used to store the Tweeter app data.
- The DC/OS [Kafka][2] publish-subscribe message service receives tweets from Cassandra and routes them to Zeppelin for real-time analytics.
- The [Marathon load balancer (Marathon-LB)][12] is an HAProxy based load balancer for Marathon only. It is useful when you require external routing or layer 7 load balancing features.
- Zeppelin is an interactive analytics notebook that works with DC/OS Spark on the backend to enable interactive analytics and visualization. Because it's possible for Spark and Zeppelin to consume all of your cluster resources, you must specify a maximum number of cores for the Zeppelin service.

This tutorial demonstrates how you can build a complete IoT pipeline on DC/OS in about 15 minutes! You will learn:

*   [How to install DC/OS services.](#install-services)
*   [How to add apps to DC/OS Marathon.](#deploy-the-containerized-app)
*   [How to route apps](#deploy-the-containerized-app) to the public node with the [Marathon load balancer][5].
* [How to add streaming analytics.](#add-streaming-analytics)

**Prerequisites:**

*  [DC/OS](https://dcos.io/docs/1.9/administration/installing/) installed with at least 5 [private agents][6] and 1 [public agent][6].
*  [DC/OS CLI](https://dcos.io/docs/1.9/usage/cli/install/) installed.
*  The public IP address of your public agent node. After you have installed DC/OS with a public agent node declared, you can [navigate to the public IP address][9] of your public agent node.

# Install services
From the DC/OS web interface [**Universe**](https://dcos.io/docs/1.9/usage/webinterface/#universe) tab, install Cassandra, Kafka, Marathon-LB, and Zeppelin.

__Tip:__ You can also install DC/OS packages from the DC/OS CLI with the [`dcos package install`][11] command.

1.  Find the **cassandra** package and click the **Install Package** button and accept the default installation. Cassandra will spin up to at least 3 nodes.
1.  Find the **kafka** package and click the **Install Package** button and accept the default installation. Kafka will spin up 3 brokers.
1.  Find the **marathon-lb** package and click the **Install Package** button and accept the default installation.
1.  Install Zeppelin.
    1.  Find the **zeppelin** package and click the **Install Package** button and choose the **Advanced Installation** option.
    1.  Click the **spark** tab and set `cores_max` to `8`.
    1.  Click **Review and Install** and complete your installation.
1.  Monitor the **Services** tab to watch as your microservices are deployed on DC/OS. You will see the Health status go from Idle to Unhealthy, and finally to Healthy as the nodes come online. This may take several minutes.

    **Tip:** It can take up to 10 minutes for Cassandra to initialize with DC/OS because of race conditions.

    ![Deployed services](./img/tweeter-deployed-services.png)

# Deploy the containerized app

In this step you deploy the containerized Tweeter app to a public node.

1.  Clone the [Tweeter][13] GitHub repository to your local directory.

    ```bash
    $ git clone git@github.com:mesosphere/tweeter.git
    ```

2.  Add the `HAPROXY_0_VHOST` label to the `tweeter.json` Marathon app definition file. `HAPROXY_0_VHOST` exposes Nginx on the external load balancer with a virtual host. The `HAPROXY_0_VHOST` value is the IP or hostname of your [public agent][9] node, or may be an external loadbalancer IP or hostname depending on your cluster configuration. 

    ```json
      ],
      "labels": {
        "HAPROXY_GROUP": "external",
        "HAPROXY_0_VHOST": "<Master-Public-IP>"
      }
    }
    ```

    For example, if you are using AWS, this can be your public ELB hostname. 

    **Important:** You must remove the leading `http://` and the trailing `/`.
    
    It should look similar to this:

    ```bash
      ],
      "labels": {
        "HAPROXY_GROUP": "external",
        "HAPROXY_0_VHOST": "joel-oss-publicsl-e21skwtlxt0c-2029962837.us-west-2.elb.amazonaws.com"
      }
    }
    ```

4.  Install and deploy Tweeter with this command.

    ```bash
    $ dcos marathon app add tweeter.json
    ```

    **Tip:** The `instances` parameter in `tweeter.json` specifies the number of app instances. Use the following command to scale your app up or down:

    ```bash
    $ dcos marathon app update tweeter instances=<number_of_desired_instances>
    ```

    The service talks to Cassandra via `node-0.cassandra.mesos:9042`, and Kafka via `broker-0.kafka.mesos:9557` in this example. Traffic is routed via the Marathon-LB (Marathon-LB) because you added the HAPROXY_0_VHOST tag on the `tweeter.json` definition.

1.  Go to the DC/OS web interface to verify your app is up and healthy. Then, navigate to the [public agent][9] node, or your ELB address if that's what was configured in the previous step, to see the Tweeter UI and post a Tweet.

    ![Tweeter][14]

# Post 100K Tweets

Use the `post-tweets.json` app a large number of Shakespeare tweets from a file:

        $ dcos marathon app add post-tweets.json


The app will post more than 100k tweets one by one, so you'll see them coming in steadily when you refresh the page.

The post-tweets app works by streaming to the VIP `1.1.1.1:30000`. This address is declared in the `cmd` parameter of the `post-tweets.json` app definition. The app uses the service discovery and load balancer service that is installed on every DC/OS node. You can see the Tweeter app defined with this VIP in the json definition under `VIP_0`.

# Add Streaming Analytics

Next, you'll perform real-time analytics on the stream of tweets coming in from Kafka.

1.  Navigate to Zeppelin at `https://<master_ip>/service/zeppelin/`, click **Import Note** and import `tweeter-analytics.json`. Zeppelin is preconfigured to execute Spark jobs on the DC/OS cluster, so there is no further configuration or setup required. Be sure to use `https://`, not `http://`.

    **Tip:** Your master IP address is the URL of the DC/OS web interface.

2.  Navigate to **Notebook** > **Tweeter Analytics**.

3.  Run the Load Dependencies step to load the required libraries into Zeppelin.

4.  Run the Spark Streaming step, which reads the tweet stream from ZooKeeper and puts them into a temporary table that can be queried using SparkSQL.

5.  Run the Top Tweeters SQL query, which counts the number of tweets per user using the table created in the previous step. The table updates continuously as new tweets come in, so re-running the query will produce a different result every time.

![Top Tweeters][16]

 [1]: https://github.com/dcos/examples/tree/master/1.8/cassandra
 [2]: https://github.com/dcos/examples/tree/master/1.8/kafka
 [3]: https://github.com/dcos/examples/tree/master/1.8/spark
 [4]: https://github.com/dcos/examples/tree/master/1.8/zeppelin
 [5]: https://github.com/mesosphere/marathon-lb
 [6]: https://dcos.io/docs/1.9/overview/concepts/
 [7]: https://dcos.io/docs/1.9/administration/installing/cloud/
 [8]: https://dcos.io/docs/1.9/administration/installing/custom/
 [9]: https://dcos.io/docs/1.9/administration/locate-public-agent/
 [10]: ./img/webui-universe-install.png
 [11]: https://dcos.io/docs/1.9/usage/cli/command-reference/
 [12]: https://dcos.io/docs/1.9/usage/service-discovery/marathon-lb/
 [13]: https://github.com/mesosphere/tweeter
 [14]: ./img/tweeter.png
 [16]: ./img/top-tweeters.png
