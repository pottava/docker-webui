Docker Web-UI
---

## Description

A web user-interface for [docker](https://www.docker.com/).  
([日本語はこちら](https://github.com/pottava/docker-webui/blob/master/README-ja.md))

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
