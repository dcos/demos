# Fast Data: Financial Transaction Processing


## Local development

### Preparation

Using [DC/OS tunneling](https://dcos.io/docs/1.8/administration/access-node/tunnel/): 

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
$ ./fintrans
2016/12/05 11:50:27 > message sent to partition 0 at offset 0
```

Then, in order to consume the messages:

```bash
$ dcos node ssh --master-proxy --leader
...
core@ip-10-0-6-69 ~ $ docker run -it mesosphere/kafka-client
...
root@e7c989566a22:/bin# ./kafka-console-consumer.sh --zookeeper leader.mesos:2181/dcos-service-kafka --topic fintrans --from-beginning
123
```

As a result you should see `123` here.

