FROM golang:1.3.3

RUN mkdir /srv/discovery

EXPOSE 8087

ADD . /srv/discovery
RUN cd /srv/discovery && ./build
CMD cd /srv/discovery && ./devweb
