FROM alpine:3.6

LABEL com.github.pottava.application="docker-webui" \
      com.github.pottava.description="A web user-interface for docker." \
      com.github.pottava.usage="docker run --rm -p 80:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui" \
      com.github.pottava.license="MIT"

ENV APP_MODE=production \
    APP_PORT=9000 \
    GOPATH=/go

RUN apk add --no-cache ca-certificates
RUN apk --no-cache add --virtual build-dependencies bash gcc musl-dev openssl go git \

    # Install go 1.8
    && GOLANG_VERSION=1.8.3 \
    && GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    && GOLANG_SRC_SHA256=5f5dea2447e7dcfdc50fa6b94c512e58bfba5673c039259fd843f68829d99fa6 \
    && export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
    && wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
    && echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && wget -q https://raw.githubusercontent.com/docker-library/golang/master/1.8/alpine3.6/no-pic.patch \
    && cd /usr/local/go/src \
    && patch -p2 -i /no-pic.patch \
    && ./make.bash \
    && mkdir -p /go/src /go/bin \
    && chmod -R 777 /go \

    # Compile docker-webui 
    && go get -u github.com/pottava/docker-webui \
    && mv /go/bin/docker-webui /usr/bin \

    # Clean up
    && apk del --purge -r build-dependencies \
    && rm -rf /usr/local/go /usr/lib/go /golang.tar.gz /*.patch /go/pkg /go/bin \
        /go/src/golang.org \
        /go/src/github.com/[^p]* \
        /go/src/github.com/pottava/docker-webui/.git* \
        /go/src/github.com/pottava/docker-webui/[^a]* \
        /go/src/github.com/pottava/docker-webui/app/[^av]* \
        /go/src/github.com/pottava/docker-webui/app/assets/scss \
        /go/src/github.com/pottava/docker-webui/app/assets/js/clients \
        /go/src/github.com/pottava/docker-webui/app/assets/js/containers \
        /go/src/github.com/pottava/docker-webui/app/assets/js/images

VOLUME ["/var/run/docker.sock"]
EXPOSE 9000

ENTRYPOINT ["docker-webui"]
