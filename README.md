# Auto Requester

一定範囲の間隔で指定の URL へリクエストを自動送信します。<br>
Go で実装されたシンプルなライブラリです。<br>
docker image を利用することですぐに利用可能です。<br>

## Usage

1. compose.yml に auto-requester を追加します
2. `$ docker compose up auto-requester`
3. コンテナが起動すると、environment で指定した URL に対しリクエストの送信を始めます 🥳

### Stop the job

以下のコマンドを使用して、`http://localhost:8080/stop`に GET リクエストを送信し、ジョブを停止します。

```sh
curl http://localhost:8080/stop
```

### Restart the job

以下のコマンドを使用して、`http://localhost:8080/start`に GET リクエストを送信し、ジョブを再開します。

```sh
curl http://localhost:8080/start
```

### リクエストボディの設定

Volume mount を用いて`/etc/app/body.json`に送信したいリクエストボディを設定できます。<br>
GET, DELETE リクエストの場合は送信されません。

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
      - '8080:8080'
    environment:
      - PORT=8080
      - FILE_PATH=/etc/app/body.json
      - INTERVAL_MIN_SEC=4
      - INTERVAL_MAX_SEC=6
      - TARGET_URL=https://httpbin.org/post
      - HTTP_METHOD=POST
      - CONTENT_TYPE=application/json
      - RANDOMIZE=true
```

### environments

| env              | description                                                                                                                                                                     | default               |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| PORT             | ジョブの停止/再開を受け付けるサーバーのポート番号<br>ports で指定したコンテナ側のポートと一致させてください                                                                     | 8080                  |
| FILE_PATH        | リクエストボディとなるファイルのパス<br>マウント先のパスと一致させてください                                                                                                    | /etc/app/body.json    |
| INTERVAL_MIN_SEC | リクエストを送信する最小の間隔(秒)                                                                                                                                              | 3                     |
| INTERVAL_MAX_SEC | リクエストを送信する最大の時間間隔(秒)<br>INTERVAL_MIN_SEC との間のランダムな間隔でリクエストが実行されます。<br> 一定間隔にしたい場合、INTERVAL_MIN_SEC と同じ値にしてください | 5                     |
| TARGET_URL       | リクエストの送信先 URL                                                                                                                                                          | http://localhost:3000 |
| HTTP_METHOD      | リクエストの HTTP メソッド                                                                                                                                                      | GET                   |
| CONTENT_TYPE     | リクエストの Content-Type ヘッダー                                                                                                                                              | application/json      |
| RANDOMIZE        | true にするとリクエストボディの json が配列の場合、配列内の要素から 1 つを毎回ランダムに選択します                                                                              | true                  |

## Link

- docker hub
  - [https://hub.docker.com/r/ishikawa096/auto-requester](https://hub.docker.com/r/ishikawa096/auto-requester)

## 📝TODO

- JSON 以外のリクエストボディに対応
- リクエストに認証ヘッダーを追加可能にする

<details>
<summary>✅Completed</summary>

- 色付きログ
- RANDOM_BODY の実装
- 自動テスト
- 異常系のテストコード追加
- docker hub へ自動デプロイ
- http リクエストを送信する
- 間隔、リクエストボディを指定可能にする
- スタート/ストップ操作のための API エンドポイント
- Docker image をアップロード

</details>
