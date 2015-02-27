FROM scratch
MAINTAINER "CoreOS, Inc"
EXPOSE 8087

# You need to build bin/discovery-linux64-static first; check build-static.

ADD bin/discovery-linux64-static /discovery

CMD ["--addr=:8087"]
ENTRYPOINT ["/discovery"]
