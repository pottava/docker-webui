# Dockerfile

FROM pottava/golang
MAINTAINER @pottava

LABEL com.github.pottava.application="docker-webui" \
      com.github.pottava.description="A web user-interface for docker." \
      com.github.pottava.usage="docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui" \
      com.github.pottava.license="MIT"

RUN go get -u github.com/fsouza/go-dockerclient
RUN go get -u golang.org/x/net/context
