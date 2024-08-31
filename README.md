# Auto Requester

一定範囲の間隔で指定のURLへリクエストを自動送信します。<br>
Goで実装されたシンプルなライブラリです。<br>
docker imageを利用することですぐに利用可能です。<br>

## Usage

1. compose.ymlにauto-requesterを追加します
2. `$ docker compose up auto-requester`
3. コンテナが起動すると、environmentで指定したURLに対しリクエストの送信を始めます🥳

### サーバーの停止

以下のコマンドを使用して、`http://localhost:8080/stop`にGETリクエストを送信し、サーバーを停止します。

```sh
curl http://localhost:8080/stop
```

### サーバーの再開

以下のコマンドを使用して、`http://localhost:8080/start`にGETリクエストを送信し、サーバーを再開します。

```sh
curl http://localhost:8080/start
```

### リクエストボディの設定

Volume mountを用いて`/etc/app/body.json`に送信したいリクエストボディを設定できます。<br>
GETリクエストの場合は送信されません。

```yml
    volumes:
      - ./path/to/your_body.json:/etc/app/body.json
```

### compose.yml example

```yml
services:
  auto-requester:
    image: ishikawa096/auto-requester:latest
    volumes:
      - ./body.json:/etc/app/body.json
    ports:
      - "8080:8080"
    environment:
      - INTERVAL_MIN_SEC=4
      - INTERVAL_MAX_SEC=6
      - TARGET_URL=https://httpbin.org/post
      - HTTP_METHOD=POST
      - CONTENT_TYPE=application/json
      - RANDOM_BODY=true
```

### environments

| env         | description                                             | default |
|-------------------|--------------------------------------------------|-------------|
| INTERVAL_MIN_SEC  | リクエストを送信する最小の間隔(秒)           | 3           |
| INTERVAL_MAX_SEC  | リクエストを送信する最大の時間間隔(秒)<br>INTERVAL_MIN_SECとの間のランダムな間隔でリクエストが実行されます。<br>一定間隔にしたい場合、INTERVAL_MIN_SECと同じ値にしてください            | 5           |
| TARGET_URL        | リクエストの送信先URL                 | http://localhost:3000 |
| HTTP_METHOD       | リクエストのHTTPメソッド。<br>使用可能: GET, POST                         | GET        |
| CONTENT_TYPE      | リクエストに含めるHTTPヘッダー                   | application/json |
| (TODO)RANDOM_BODY       | リクエストボディのjsonが配列の場合、配列内の要素から1つを毎回ランダムに選択し、リクエストボディとして用います| true        |

## Link

- docker hub
  - [https://hub.docker.com/r/ishikawa096/auto-requester](https://hub.docker.com/r/ishikawa096/auto-requester)

## TODO

- RANDOM_BODYの実装
- リクエストに認証ヘッダーを追加可能にする
- 異常系のテストコード追加
- 自動テスト
- docker hubへ自動デプロイ
