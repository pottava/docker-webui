Usage with VirtualBox, CoreOS & Docker containers
---

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

[hhttp://docker-webui.localhost.com/](http://docker-webui.localhost.com/)

### 6. Test the application

```shell
$ vagrant ssh -c "docker run --rm -v /home/core/share:/go/src/github.com/pottava/docker-webui pottava/docker-webui:base go test github.com/pottava/docker-webui/..."
```

### 7. Teardown the VM

```shell
$ vagrant halt
```
