Docker Web-UI
---

## 概要

[docker](https://www.docker.com/)を Webから操作するためのツールです。

## できること

### Dockerコンテナを操作する

* コンテナの検索（状態、検索文字列によるフィルタリング）
* inspect, top, stats, logs, diff, rename, commitコマンドの実行
* コンテナの start, stop, restart, rmの実行
* 起動元になっているイメージの検索

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/containers.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### コンテナごとのログをみる

* 指定時間秒ごとに最新の 10〜200行を表示します
* 現在の CPUや メモリの状況も確認できます

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### 起動しているコンテナの stats一覧

* 起動中コンテナごとの CPUと メモリの状況をグラフで表示します
* データ取得時間など詳細データは下部にテーブルで表示

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/stats.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### Dockerイメージを操作する

* 新規 Dockerイメージの pull
* 同名のイメージを改めて pull
* イメージの検索
* inspect, history, tag, rmiコマンドの実行

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/images.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

## 使い方

### 1. アプリケーションを起動します

Dockerコンテナとしてなら

`$ docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui`

Go言語環境のある方はこちらでも

`$ go get github.com/pottava/docker-webui`  
`$ docker-webui`

### 2. 以下のURLを開きます

[http://localhost:9000/](http://localhost:9000/)
