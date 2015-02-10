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

## Run in docker

```
docker run -d -p 8087:8087 <#Image Name> 
```

## ENV list
there are two ENVs to control discovery:

1. `DISCOVERY_ROOT_URL` : default https://discovery.etcd.io
2. `DISCOVERY_ORIGIN_ADDR` : default http://127.0.0.1:4001
3. `DISCOVERY_INIT_LEADER` : default 127.0.0.1:4001
