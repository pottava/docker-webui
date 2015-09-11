# Dockerfile

FROM golang:1.5.1
MAINTAINER @pottava

RUN go get github.com/pottava/docker-webui

VOLUME ["/var/run/docker.sock"]
EXPOSE 9000

ENTRYPOINT ["/go/bin/docker-webui"]
