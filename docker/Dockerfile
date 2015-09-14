# Dockerfile

FROM golang:1.5.1
MAINTAINER @pottava

LABEL com.github.pottava.application="docker-webui" \
      com.github.pottava.description="A web user-interface for docker." \
      com.github.pottava.usage="docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui" \
      com.github.pottava.license="MIT"

RUN go get github.com/pottava/docker-webui

VOLUME ["/var/run/docker.sock"]
EXPOSE 9000

ENTRYPOINT ["/go/bin/docker-webui"]
