# Fast Data: Financial Transaction Processing


## Local development

We're using Shopify's [sarama](https://godoc.org/github.com/Shopify/sarama) package to communicate with Kafka.

### Preparation

[Install](https://github.com/dcos/examples/tree/master/1.8/kafka) the Apache Kafka package and then figure out where the brokers are:

```bash
$ dcos kafka connection

{
  "address": [
    "10.0.3.178:9398",
    "10.0.3.176:9382",
    "10.0.3.177:9363"
  ],
  "zookeeper": "master.mesos:2181/dcos-service-kafka",
  "dns": [
    "broker-0.kafka.mesos:9398",
    "broker-1.kafka.mesos:9382",
    "broker-2.kafka.mesos:9363"
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
$ ./fintrans --broker broker-1.kafka.mesos:9382
INFO[0000] &sarama.ProducerMessage{Topic:"Tokyo", Key:sarama.Encoder(nil), Value:"226 793 864262", Metadata:interface {}(nil), Offset:7, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0003] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"543 27 655474", Metadata:interface {}(nil), Offset:0, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0005] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"36 610 702712", Metadata:interface {}(nil), Offset:1, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0007] &sarama.ProducerMessage{Topic:"NYC", Key:sarama.Encoder(nil), Value:"279 526 765420", Metadata:interface {}(nil), Offset:2, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
INFO[0009] &sarama.ProducerMessage{Topic:"London", Key:sarama.Encoder(nil), Value:"456 698 575718", Metadata:interface {}(nil), Offset:3, Partition:0, Timestamp:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, retries:0, flags:0}
...
```

Then, in order to consume the messages:

```bash
$ dcos node ssh --master-proxy --leader
...
core@ip-10-0-6-69 ~ $ docker run -it mesosphere/kafka-client
...
root@e7c989566a22:/bin# ./kafka-console-consumer.sh --zookeeper leader.mesos:2181/dcos-service-kafka --topic NYC --from-beginning
543 27 655474
36 610 702712
279 526 765420
445 284 174140
870 826 914135
26 23 351148
64 751 531640
922 354 673924
763 347 600496
^CProcessed a total of 9 messages
```

As a result you should see something like above until you hit `CTRL+C`. The wire format of the messages is: `source_account target_account amount`.

