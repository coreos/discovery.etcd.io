# etcd-discovery

etcd-discovery service powers the public service at https://discovery.etcd.io. The API is documented in

https://github.com/coreos/etcd/tree/master/Documentation/cluster-discovery.md
https://github.com/coreos/etcd/tree/master/Documentation/discovery-protocol.md

## Development

etcd-discovery uses devweb for easy development. It is simple to get started:

```
./devweb
curl --verbose -X PUT localhost:8087/new
```
