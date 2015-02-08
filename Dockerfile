FROM google/golang
MAINTAINER Meng Zhuo "mengzhuo1203@gmail.com"

RUN mkdir /srv/discovery

EXPOSE 8087

ADD . /srv/discovery
RUN cd /srv/discovery && ./build

CMD cd /srv/discovery && ./bin/discovery.etcd.io
