# discovery.etcd.io

This code powers the public service at https://discovery.etcd.io. The API is
documented in the etcd clustering documentation:

https://github.com/coreos/etcd/blob/master/Documentation/clustering.md#public-etcd-discovery-service

## Docker Container

You may run the service in a docker container:

```
docker pull quay.io/coreos/discovery.etcd.io
docker run -d -p 80:8087 quay.io/coreos/discovery.etcd.io
```

## Development

discovery.etcd.io uses devweb for easy development. It is simple to get started:

```
./devweb
curl --verbose -X PUT localhost:8087/new
```

## Environment Variables 

You can set the discovery service with your address by `BASE_URL`. When you visit `http://<your discovery service>/new`, server will return `http://<BASE_URL>/<token>`  

If your discovery service run as a container , maybe you need `ETCD_CONN`. You can set the etcd connection like `http://172.17.42.1:2379`, then discovery service access etcd from the docker bridge.

Both of these are Optional. If you do not set these environment variables, no difference between previous.
### Example:
```
docker run -d -p 80:8087 -e ETCD_COON=http://172.17.42.1:2379 -e BASE_URL=http://192.168.2.254 <image id>
```