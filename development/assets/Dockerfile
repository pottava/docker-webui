FROM pottava/gulp:3.9-docker

ADD . /gulp
WORKDIR /gulp

RUN apk add --no-cache make g++ python-dev \
    && npm install
