Docker Web-UI
---

## Description

A web user-interface for [docker](https://www.docker.com/).

## Basic Usage

### 1. Run as a docker container

```shell
$ docker run --rm -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui:latest
```

### 2. Access the application

[http://localhost:9000/](http://localhost:9000/)

## Usage with Golang

### 1. Install go binary

```shell
$ go get github.com/pottava/docker-webui
```

### 2. Run this application

```shell
$ docker-webui
```

### 3. Access the application

[http://localhost:9000/](http://localhost:9000/)

## Usage with VirtualBox, CoreOS & Docker containers

### 1. Install VirtualBox & Vagrant

- [VirtualBox](https://www.virtualbox.org/)
- [Vagrant](http://www.vagrantup.com/)

### 2. Install vagrant-hostsupdater plugin

```shell
$ vagrant plugin install vagrant-hostsupdater
```

### 3. Change your working directory to vagrant folder

```shell
$ cd /path/to/this-repository-root/vagrant
```

### 4. Create a virtual machine

```shell
$ vagrant up
```

### 5. Confirm whether a service is running

[hhttp://docker-webui.local/](http://docker-webui.local/)

### 6. Test the application

```shell
$ vagrant ssh -c "docker run --rm -v /home/core/share:/go/src/github.com/pottava/docker-webui pottava/docker-webui:base go test github.com/pottava/docker-webui/..."
```

### 7. Teardown the VM

```shell
$ vagrant halt
```

## Contribution

1. Fork ([https://github.com/pottava/docker-webui/fork](https://github.com/pottava/docker-webui/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request

## Copyright and license

Code and documentation copyright 2015 SUPINF Inc. Code released under the [MIT license](https://github.com/pottava/docker-webui/blob/master/LICENSE).
