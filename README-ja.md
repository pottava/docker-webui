Docker Web-UI
---

## 概要

[docker](https://www.docker.com/)を Webから操作するためのツールです。

## 使い方

### 1. アプリケーションを起動します

Dockerコンテナとしてなら

`$ docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui`

Go言語環境のある方はこちらでも

`$ go get github.com/pottava/docker-webui`  
`$ docker-webui`

### 2. 以下のURLを開きます

[http://localhost:9000/](http://localhost:9000/)
