# discovery.etcd.io

This code powers the public service at https://discovery.etcd.io. The API is
documented in the [etcd clustering documentation](https://github.com/coreos/etcd/blob/master/Documentation/dev-internal/discovery_protocol.md#public-discovery-service).

# Configuration

The service has three configuration options, and can be configured with either
runtime arguments or environment variables.

* `--addr` / `DISC_ADDR`: the address to run the service on, including port.
* `--host` / `DISC_HOST`: the host url to prepend to `/new` requests.
* `--etcd` / `DISC_ETCD`: the url of the etcd endpoint backing the instance.

## Docker Container

You may run the service in a docker container:

```
docker pull quay.io/coreos/discovery.etcd.io
docker run -d -p 80:8087 -e DISC_ETCD=http://etcd.example.com:2379 -e DISC_HOST=http://discovery.example.com quay.io/coreos/discovery.etcd.io
```

## Development

discovery.etcd.io uses devweb for easy development. It is simple to get started:

```
./devweb
curl --verbose -X PUT localhost:8087/new
```
