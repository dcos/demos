# Fast Data: Financial Transaction Processing


## Local development

We're using Shopify's [sarama](https://godoc.org/github.com/Shopify/sarama) package to communicate with Kafka.

### Preparation

[Install](https://github.com/dcos/examples/tree/master/1.8/kafka) the Apache Kafka package with the following [options](kafka-config.json):

```bash
$ dcos package install kafka --options=kafka-config.json
```

Next, figure out where the broker is:

```bash
$ dcos kafka connection

{
  "address": [
    "10.0.3.178:9398"
  ],
  "zookeeper": "master.mesos:2181/dcos-service-kafka",
  "dns": [
    "broker-0.kafka.mesos:9398"
  ],
  "vip": "broker.kafka.l4lb.thisdcos.directory:9092"
}
```

Now, using [DC/OS tunneling](https://dcos.io/docs/1.8/administration/access-node/tunnel/): 

```bash
$ sudo dcos tunnel vpn --client=/Applications/Tunnelblick.app/Contents/Resources/openvpn/openvpn-2.3.12/openvpn
Password:
*** Unknown ssh-rsa host key for 35.156.70.254: 13ec7cde1d3967d2371eb375f48c4690

ATTENTION: IF DNS DOESN'T WORK, add these DNS servers!
198.51.100.1
198.51.100.2
198.51.100.3

Waiting for VPN server in container 'openvpn-6nps1efm' to come up...

VPN server output at /tmp/tmpn34d7n0d
VPN client output at /tmp/tmpw6aq3v4z
```

Note that it may be necessary to [add the announced DNS servers]( https://support.apple.com/kb/PH18499?locale=en_US) as told by Tunnelblick, and make sure they are the first in the list.

### Producing and consuming messages

With the VPN tunnel enabled, we can run the fintrans generator:

```bash
$ ./fintrans --broker broker-0.kafka.mesos:9398
INFO[0001] &sarama.ProducerMessage{Topic:"London", Key:sarama.Encoder(nil), Value:"678 816 2957", Metadata:interface {}(nil), Offset:10, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0003] &sarama.ProducerMessage{Topic:"SF", Key:sarama.Encoder(nil), Value:"762 543 6395", Metadata:interface {}(nil), Offset:4, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0005] &sarama.ProducerMessage{Topic:"London", Key:sarama.Encoder(nil), Value:"680 840 8115", Metadata:interface {}(nil), Offset:11, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0007] &sarama.ProducerMessage{Topic:"SF", Key:sarama.Encoder(nil), Value:"363 101 9878", Metadata:interface {}(nil), Offset:5, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0009] &sarama.ProducerMessage{Topic:"SF", Key:sarama.Encoder(nil), Value:"302 505 5777", Metadata:interface {}(nil), Offset:6, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0011] &sarama.ProducerMessage{Topic:"London", Key:sarama.Encoder(nil), Value:"848 948 2683", Metadata:interface {}(nil), Offset:12, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0013] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"611 695 5484", Metadata:interface {}(nil), Offset:9, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0015] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"396 465 6789", Metadata:interface {}(nil), Offset:10, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0018] &sarama.ProducerMessage{Topic:"Moscow", Key:sarama.Encoder(nil), Value:"132 570 3197", Metadata:interface {}(nil), Offset:9, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0020] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"607 672 9732", Metadata:interface {}(nil), Offset:11, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
^C
```

Then, in order to consume the messages:

```bash
$ dcos node ssh --master-proxy --leader
...
core@ip-10-0-6-69 ~ $ docker run -it mesosphere/kafka-client
...
root@e7c989566a22:/bin# ./kafka-console-consumer.sh --zookeeper leader.mesos:2181/dcos-service-kafka --topic NYC --from-beginning
611 695 5484
396 465 6789
607 672 9732
^CProcessed a total of 3 messages
```

As a result you should see something like above until you hit `CTRL+C`. The wire format of the messages is: `source_account target_account amount`.

Note: if you want to reset the topics, do a `dcos kafka topic list` and `dcos kafka topic delete XXX` with `XXX` being one of the listed topics.
