Docker Web-UI
---

## 概要

[docker](https://www.docker.com/)を Webから操作するためのツールです。

## できること

### Dockerコンテナを操作する

* コンテナの検索（状態、ラベル、検索文字列によるフィルタリング）
* inspect, top, stats, logs, diff, rename, commitコマンドの実行
* コンテナの start, stop, restart, kill, rmの実行
* 起動元になっているイメージの検索

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/containers.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### 起動しているコンテナの ログ一覧

* ラベルによるログ表示対象コンテナ絞り込み
* 指定時間秒ごとに最新の 10〜200行を表示します

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### 起動しているコンテナの stats一覧

* ラベルによるスタッツ表示対象コンテナ絞り込み
* 起動中コンテナごとの CPUと メモリの状況をグラフで表示します
* データ取得時間など詳細データは下部テーブルに表示

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/stats.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### コンテナごとのログ & スタッツをみる

* 指定時間秒ごとに最新の 10〜200行を表示します
* 現在の CPUや メモリの状況も確認できます

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/specified-logs.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

### Dockerイメージを操作する

* イメージの検索（ラベル、検索文字列によるフィルタリング）
* 新規 Dockerイメージの pull
* 同名のイメージを改めて pull
* inspect, history, tag, rmiコマンドの実行

<img alt="" src="https://raw.github.com/wiki/pottava/docker-webui/images/images.png"
  style="width: 500px;-webkit-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         -moz-box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);
         box-shadow: 4px 6px 10px 0px rgba(0,0,0,0.5);">

## 使い方

### 1. アプリケーションを起動します

docker-composeのサービスとして起動するなら

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

シンプルにDockerコンテナとしてなら

`$ docker run -p 9000:9000 --rm -v /var/run/docker.sock:/var/run/docker.sock pottava/docker-webui`  
`$ docker run -p 9000:9000 --rm -e DOCKER_HOST -e DOCKER_CERT_PATH=/etc/docker-compose/cert -v $DOCKER_CERT_PATH:/etc/docker-compose/cert pottava/docker-webui`

Go言語環境のある方はこちらでも

`$ go get github.com/pottava/docker-webui`  
`$ docker-webui`

### 2. 以下のURLを開きます

[http://localhost:9000/](http://localhost:9000/)


## 起動時オプション

環境変数、または [config.json](https://github.com/pottava/docker-webui/blob/master/config.json) を使って起動時に設定変更できます。

Option (config.json)   | Environment Variables     | Description                                       | Default
---------------------- | ------------------------- | ------------------------------------------------- | ---------
Name                   | APP_NAME                  | アプリケーションの名前                                 | 'docker web-ui'
Port                   | APP_PORT                  | アプリが利用するポート                                 | 9000
　 | CONFIG_FILE_PATH | config.json の配置パス | '/etc/docker-webui/config.json'
ViewOnly               | APP_VIEW_ONLY             | Dockerの状態変更系アクションを抑制します                 | false
LogLevel               | APP_LOG_LEVEL             | 1:fatal, 2:err, 3:warn, 4:info, 5:debug, 6:trace  | 4
LabelOverrideNames     | APP_LABEL_OVERRIDE_NAMES  | コンテナの表示名を特定のラベルの値に上書きできます           |
LabelFilters           | APP_LABEL_FILTERS         | フィルタリングに利用するラベルを指定できます               | ['all']
DockerEndpoints        | DOCKER_HOST               | docker APIのエンドポイント (tcp or socket)           | [`unix:///var/run/docker.sock`]
DockerCertPath         | DOCKER_CERT_PATH          | TLS接続に使う証明書があれば、その絶対パス                 | ['']
DockerPullBeginTimeout | DOCKER_PULL_BEGIN_TIMEOUT | docker pull開始までのタイムアウト時間                   | 3 * time.Minute
DockerPullTimeout      | DOCKER_PULL_TIMEOUT       | docker pull実行のタイムアウト時間                      | 2 * time.Hour
DockerStatTimeout      | DOCKER_STAT_TIMEOUT       | docker stat実行のタイムアウト時間                      | 5 * time.Second
DockerStartTimeout     | DOCKER_START_TIMEOUT      | docker start実行のタイムアウト時間                     | 10 * time.Second
DockerStopTimeout      | DOCKER_STOP_TIMEOUT       | docker stop実行のタイムアウト時間                      | 10 * time.Second
DockerRestartTimeout   | DOCKER_RESTART_TIMEOUT    | docker restart実行のタイムアウト時間                   | 10 * time.Second
DockerKillTimeout      | DOCKER_KILL_TIMEOUT       | docker kill実行のタイムアウト時間                      | 10 * time.Second
DockerRmTimeout        | DOCKER_RM_TIMEOUT         | docker rm実行のタイムアウト時間                        | 5 * time.Minute
DockerCommitTimeout    | DOCKER_COMMIT_TIMEOUT     | docker commit実行のタイムアウト時間                    | 30 * time.Second
StaticFileHost         | APP_STATIC_FILE_HOST      | 静的ファイル配信ホスト名                               | 
StaticFilePath         | APP_STATIC_FILE_PATH      | ホスト上の静的ファイル配置パス                          | '$GOPATH + /src/github.com/pottava/docker-webui/app'
PathPrefix             | APP_PATH_PREFIX           | パスベースルーティングなどのための URL パスプリフィックス   | 
PreventSelfStop        | APP_PREVENT_SELF_STOP     | このアプリ自身をWebUIから停止することを防ぎます            | true
HiddenContainers       | APP_HIDDEN_CONTAINERS     | 画面上表示したくないコンテナを指定できます                 | []
