Docker Web-UI
---

## Description

A web user-interface for [docker](https://www.docker.com/).  
([日本語はこちら](https://github.com/pottava/docker-webui/blob/master/README-ja.md))

## Features

### handling docker containers

* search containers (by its status, query strings)
* inspect, top, stats, logs, diff, rename, commit
* start, stop, restart, rm
* search its image

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/containers.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### tailing container logs

* monitoring logs (10~200 lines each) every specified seconds
* you can check CPUs & memorys at the same time

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### watch container statistics

* display CPU & memory on running containers by charts
* display details on the bottom table

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/stats.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### handling docker images

* docker pull a new image
* docker pull the same image again
* search images
* inspect, history, tag, rmi

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/images.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

## Usage

### 1. Run the application

as a docker container

`$ docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui`

or as a go binary

`$ go get github.com/pottava/docker-webui`  
`$ docker-webui`

### 2. Access to the following URL

[http://localhost:9000/](http://localhost:9000/)



## Contribution

1. Fork ([https://github.com/pottava/docker-webui/fork](https://github.com/pottava/docker-webui/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request

## Copyright and license

Code and documentation copyright 2015 SUPINF Inc. Code released under the [MIT license](https://github.com/pottava/docker-webui/blob/master/LICENSE).
