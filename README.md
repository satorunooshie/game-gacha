# game-gacha
## 概要

SwaggerEditor: <https://editor.swagger.io> <br>
定義ファイル: `./api-document.yaml`<br>

※ Firefoxはブラウザ仕様により上記サイトからlocalhostへ向けた通信を許可していないので動作しません
- https://bugzilla.mozilla.org/show_bug.cgi?id=1488740
- https://bugzilla.mozilla.org/show_bug.cgi?id=903966

## 事前準備
### docker-composeを利用したMySQLとRedisの準備
#### MySQL
MySQLはリレーショナルデータベースの1つです。
```
$ docker-compose up mysql
```
を実行することでローカルのDocker上にMySQLサーバが起動します。<br>
<br>
初回起動時に db/init ディレクトリ内のDDL, DMLを読み込みデータベースの初期化を行います。<br>
(DDL(DataDefinitionLanguage)とはデータベースの構造や構成を定義するためのSQL文)<br>
(DML(DataManipulationLanguage)とはデータの管理・操作を定義するためのSQL文)

#### Redis
Redisはインメモリデータベースの1つです。<br>
```
$ docker-compose up redis
```
を実行することでローカルのDocker上にMySQLサーバが起動します。

### MySQLWorkbenchの設定
MySQLへの接続設定をします。
1. MySQL Connections の + を選択
2. 以下のように接続設定を行う
    ```
    Connection Name: 任意 
    Connection Method: Standard (TCP/IP)
    Hostname: 127.0.0.1 (localhost)
    Port: 3306
    Username: root
    Password: game-gacha
    Default Schema: game_gacha
    ```

### API用のデータベースの接続情報を設定する
環境変数にデータベースの接続情報を設定します。<br>
ターミナルのセッション毎に設定したり、.bash_profileで設定を行います。

Macの場合
```
$ export MYSQL_USER=root \
    MYSQL_PASSWORD=game-gacha \
    MYSQL_HOST=127.0.0.1 \
    MYSQL_PORT=3306 \
    MYSQL_DATABASE=game_gacha
```

Windowsの場合
```
$ SET MYSQL_USER=root
$ SET MYSQL_PASSWORD=game-gacha
$ SET MYSQL_HOST=127.0.0.1
$ SET MYSQL_PORT=3306
$ SET MYSQL_DATABASE=game_gacha
```

## APIローカル起動方法
```
$ go run ./cmd/main.go
```