# gRPC Learning App

gRPC学習向けの最小構成サンプルです。

- `proto`: サービス定義
- `gen/go`: Goの生成コード
- `web/src/gen`: TypeScriptの生成コード
- `cmd/server`: Goサーバー（gRPC / gRPC-Web / 静的Web配信）

gRPCの使い方の説明は `doc/grpc-usage.md` を参照してください。

## セットアップ

```bash
make setup
make deps
make generate
make web-build
make go-build
```

まとめて実行する場合:

```bash
make build
```

## 起動

```bash
make run
```

- Web: http://localhost:8080
- gRPC-Web endpoint: `/todo.v1.TodoService/*`

`8080` が埋まっている場合はポートを変更できます。

```bash
PORT=18080 make run
```

## 学習用コマンド

### 1) gRPCで直接呼び出し（grpcurl）

```bash
grpcurl \
  -plaintext \
  -import-path proto \
  -proto todo/v1/todo.proto \
  -d '{"title":"grpcurlから追加"}' \
  localhost:8080 \
  todo.v1.TodoService/AddTodo
```

```bash
grpcurl \
  -plaintext \
  -import-path proto \
  -proto todo/v1/todo.proto \
  -d '{}' \
  localhost:8080 \
  todo.v1.TodoService/ListTodos
```

### 2) proto変更後の再生成

```bash
make generate
```
