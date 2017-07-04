Docker Web-UI
---

Supported tags and respective `Dockerfile` links:
・latest ([production/Dockerfile](https://github.com/pottava/docker-webui/blob/master/production/Dockerfile))

## Description

A web user-interface for [docker](https://www.docker.com/).  
([日本語はこちら](https://github.com/pottava/docker-webui/blob/master/README-ja.md))

## Features

### handling docker containers

* search containers (by status, labels, query strings)
* inspect, top, stats, logs, diff, rename, commit
* start, stop, restart, kill, rm
* search its image

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/containers.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### tailing containers' logs

* filter target containers by labels
* monitoring logs (10~200 lines each) every specified seconds

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### watch containers' statistics

* filter target containers by labels
* display CPU & memory on running containers by charts
* details on the bottom table

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/stats.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### watch specified container's logs & stats

* monitoring logs (10~200 lines each) every specified seconds
* you can check CPUs & memorys at the same time

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/specified-logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### handling docker images

* search images (by labels, query strings)
* docker pull a new image
* docker pull the same image again
* inspect, history, tag, rmi

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/images.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

## Usage

### 1. Run the application

as a docker-compose service

```
monit:
  image: pottava/docker-webui
  ports:
    - "9000:9000"
  volumes:
    - "${DOCKER_CERT_PATH}:/etc/docker-compose/cert"
  environment:
    - DOCKER_HOST
    - DOCKER_CERT_PATH=/etc/docker-compose/cert
    - APP_LABEL_OVERRIDE_NAMES=com.docker.compose.service
    - APP_LABEL_FILTERS=com.docker.compose.service
```

or as a simple docker container

`$ docker run -p 9000:9000 --rm -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui`  
`$ docker run -p 9000:9000 --rm -e DOCKER_HOST -e DOCKER_CERT_PATH=/etc/docker-compose/cert -v $DOCKER_CERT_PATH:/etc/docker-compose/cert pottava/docker-webui`

or as a go binary

`$ go get github.com/pottava/docker-webui`  
`$ docker-webui`

### 2. Access to the following URL

[http://localhost:9000/](http://localhost:9000/)


## Startup Options

You can set environment variables or use [config.json](https://github.com/pottava/docker-webui/blob/master/config.json) to configure the application.

Option (config.json)   | Environment Variables     | Description                                       | Default
---------------------- | ------------------------- | ------------------------------------------------- | ---------
Name                   | APP_NAME                  | name of this application                          | `docker web-ui`
Port                   | APP_PORT                  | port the app is listening on                      | 9000
                       | CONFIG_FILE_PATH          | path of config.json                               |  `/etc/docker-webui/config.json`
ViewOnly               | APP_VIEW_ONLY             | if you set true, you cannot change docker state   | false
LogLevel               | APP_LOG_LEVEL             | 1:fatal, 2:err, 3:warn, 4:info, 5:debug, 6:trace  | 4
LabelOverrideNames     | APP_LABEL_OVERRIDE_NAMES  | override containers name by its label value       | ``
LabelFilters           | APP_LABEL_FILTERS         | labels for filtering containers & images          | [`all`]
DockerEndpoints        | DOCKER_HOST               | docker API endpoints (tcp or socket)              | [`unix:///var/run/docker.sock`]
DockerCertPath         | DOCKER_CERT_PATH          | set certifications' absolute path on the host     | [``]
DockerPullBeginTimeout | DOCKER_PULL_BEGIN_TIMEOUT | timeout of docker pull to start                   | 3 * time.Minute
DockerPullTimeout      | DOCKER_PULL_TIMEOUT       | timeout of docker pull                            | 2 * time.Hour
DockerStatTimeout      | DOCKER_STAT_TIMEOUT       | timeout of docker stat                            | 5 * time.Second
DockerStartTimeout     | DOCKER_START_TIMEOUT      | timeout of docker start                           | 10 * time.Second
DockerStopTimeout      | DOCKER_STOP_TIMEOUT       | timeout of docker stop                            | 10 * time.Second
DockerRestartTimeout   | DOCKER_RESTART_TIMEOUT    | timeout of docker restart                         | 10 * time.Second
DockerKillTimeout      | DOCKER_KILL_TIMEOUT       | timeout of docker kill                            | 10 * time.Second
DockerRmTimeout        | DOCKER_RM_TIMEOUT         | timeout of docker rm                              | 5 * time.Minute
DockerCommitTimeout    | DOCKER_COMMIT_TIMEOUT     | timeout of docker commit                          | 30 * time.Second
StaticFileHost         | APP_STATIC_FILE_HOST      | host name which provides static files             | ``
StaticFilePath         | APP_STATIC_FILE_PATH      | static file path on the host                      | `gopath + /src/github.com/pottava/docker-webui/app`
PreventSelfStop        | APP_PREVENT_SELF_STOP     | prevent to stop this app itself if you set true   | true
HiddenContainers       | APP_HIDDEN_CONTAINERS     | hide specified containers if you like             | []


## Contribution

1. Fork ([https://github.com/pottava/docker-webui/fork](https://github.com/pottava/docker-webui/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request

## Copyright and license

Code released under the [MIT license](https://github.com/pottava/docker-webui/blob/master/LICENSE).
