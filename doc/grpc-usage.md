# このプロジェクトでの gRPC の使い方

このアプリは、`TodoService` を題材にして gRPC の基本フローを最小構成で学ぶためのサンプルです。

## 1. gRPCとは何か

gRPC は、アプリ間通信を「RPC (Remote Procedure Call)」として扱うためのフレームワークです。  
HTTP API を手で組み立てる代わりに、「このメソッドをこの型で呼ぶ」という契約を `.proto` で定義し、サーバー/クライアントのコードを自動生成して使います。

このプロジェクトで言うと:

- 契約: `todo.proto` の `TodoService`
- サーバー: `AddTodo`, `ListTodos` を実装
- クライアント: 生成された型安全なメソッドを呼び出し

gRPC の学習で重要なのは、次の 3 点です。

- 契約駆動: `.proto` が先にある
- 型安全: Go と TypeScript の両方に同じ契約が反映される
- 通信抽象化: HTTP 詳細より「メソッド呼び出し」に集中できる

## 2. 契約を `.proto` で定義

gRPC の入り口は [proto/todo/v1/todo.proto](../proto/todo/v1/todo.proto) です。

- `message`: リクエスト/レスポンスの型
- `service`: RPC メソッド定義
- `rpc AddTodo`, `rpc ListTodos`: 今回の最小メソッド

ここが API 契約の単一情報源です。

## 3. `.proto` からコード生成

`make generate` で以下を生成します。

- Go: `gen/go/...`
- TypeScript: `web/src/gen/...`

生成設定は [buf.gen.yaml](../buf.gen.yaml) にあり、以下のプラグインを使っています。

- `protoc-gen-go`: メッセージ型
- `protoc-gen-connect-go`: Go のハンドラ/クライアント向けコード
- `protoc-gen-es`: Web 側の型とサービス定義

## 4. サーバー側での gRPC 利用

### 実装クラス

[internal/todo/service.go](../internal/todo/service.go) で RPC を実装しています。

- `AddTodo(ctx, req)`: タイトルを受け取り TODO を作成
- `ListTodos(ctx, req)`: 一覧を返却

メソッドシグネチャは `connect.Request[T]` / `connect.Response[T]` を使っています。

### エンドポイント登録

[cmd/server/main.go](../cmd/server/main.go) で次を行います。

- `todov1connect.NewTodoServiceHandler(todoService)` でハンドラ作成
- `mux.Handle(path, handler)` で RPC パスを登録
- `h2c.NewHandler(...)` で HTTP/2 cleartext を有効化

これにより、同一プロセスで Web 配信と RPC 提供をまとめています。

## 5. ブラウザ側での gRPC 利用

[web/src/main.ts](../web/src/main.ts) で gRPC-Web クライアントを作っています。

- `createGrpcWebTransport(...)`: gRPC-Web 通信設定
- `createClient(TodoService, transport)`: サービスクライアント生成
- `client.addTodo(...)`, `client.listTodos(...)`: RPC 呼び出し

ブラウザは純粋な gRPC (HTTP/2 バイナリ) を直接扱いづらいため、ここでは gRPC-Web を使っています。

## 6. 通信フロー

1. UI で入力して送信
2. `web/src/main.ts` が `AddTodo` / `ListTodos` を呼ぶ
3. サーバーの `TodoServiceHandler` が受信
4. `internal/todo/service.go` の実装が処理
5. レスポンスを gRPC-Web でブラウザに返却

## 7. 学習時に見るポイント

- `.proto` を変更すると、Go/TS の型が同時に更新される
- サーバーとクライアントで同じ契約を共有できる
- バリデーションエラーが RPC エラーとして返る

## 8. すぐ試すコマンド

```bash
make build
PORT=18080 make run
```

gRPC 直接呼び出しは [README.md](../README.md) の `grpcurl` 例を使って確認できます。
