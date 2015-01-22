FROM golang:onbuild
EXPOSE 8087

ADD . /srv/discovery.etcd.io
WORKDIR /srv/discovery.etcd.io

CMD GOPATH="${PWD}/third_party" ./devweb -addr=":8087" github.com/coreos/discovery.etcd.io/discovery
