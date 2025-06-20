# Go Web Console — プロジェクト計画書（Part 2: 機能要件詳細）

---

## 🔹 UC1. 構造化ログの可視化

### 概要
`slog` を用いて構造化されたログを出力し、レベルや期間などでフィルタ可能に。JSON形式で保存されたログを、UIまたはAPI経由で閲覧。

### API仕様
- `GET /api/logs/{level}`: ログレベル別取得（例: `info`, `warn`, `error`）
- クエリ例: `/api/logs/info?since=2024-01-01T00:00:00Z&keyword=server`

### ユースケース構造
- `GetLogsUsecase`（InputPort）
- `LogRepository`（OutputPort）
- `LogEntry`（ドメインエンティティ）

### 補足
- ログはメモリキャッシュ or ローテートJSONファイルで構成
- `slog.Handler` は `slog.NewJSONHandler` を使用
- クライアント表示は `/logs` でテンプレート描画

---

## 🔹 UC2. HTTPルーティングとパスパラメータの検証

### 概要
`ServeMux` の `{id}` や HTTPメソッド指定パターンを使用し、ルーティング挙動を可視化。複数パターン競合の処理も確認可能。

### API仕様
- `GET /api/users/{id}`
- `POST /api/tasks/{id}/restart`
- `GET /api/routes`: 登録ルート一覧の表示

### ユースケース構造
- `GetUserByIDUsecase`, `RestartTaskUsecase`
- `RouteRegistryService`（ドメインサービス）

### 補足
- `http.Request.PathValue("id")` により `id` を取得
- ServeMux登録順・優先順位の挙動比較機能あり
- GODEBUG切替（旧動作）対応：`httpmuxgo121=1`

---

## 🔹 UC3. Contextキャンセル・タイムアウト可視化

### 概要
`context` の締切／キャンセル／理由の伝播を確認するAPIを提供し、HTTPリクエストに連動するリクエストライフサイクルを体感。

### API仕様
- `GET /api/cancel?timeout=3s&bg=true`
  - 指定秒数でキャンセルされる処理を開始
  - `bg=true` の場合 `WithoutCancel` によりバックグラウンド継続

### ユースケース構造
- `CancelableTaskUsecase`
- `ContextTrace`（Value Object）
- `Logger` 出力に `Cause()` の値を含めて記録

### 補足
- `AfterFunc` による終了処理ロギング
- ゴルーチンの残存を監視（実装レベル）

---

## 🔹 UC4. Cookie / ヘッダ検証ビューア

### 概要
HTTPクッキーやヘッダの取り扱いを表示し、ブラウザの振る舞いやGoのAPI挙動を検証。複数クッキー、パーティション属性も含む。

### API仕様
- `GET /headers`: ヘッダ一覧出力
- `GET /cookies`: 現在のCookie表示
- `POST /cookies/set`: Cookie設定用

### ユースケース構造
- `GetCookiesUsecase`
- `CookieParserService`（ドメインサービス）

### 補足
- `Request.CookiesNamed("id")`
- `Cookie.Partitioned`, `Quoted` を明示的に表示
- ヘッダ解析には `ParseSetCookie`, `ParseCookie` を使用

---

## 🔹 UC5. テンプレート描画 & JS埋め込みチェック

### 概要
テンプレート内での `range` や JavaScript テンプレリテラル埋め込みの対応状況を視覚的に確認できるページ。

### ページ
- `GET /template-test`

### 内容
- `{{range 5}}●{{end}}` のループ描画
- `<script>const msg = `{{.Message}}`;</script>` 形式の検証
- JSリテラル中でのテンプレ展開制御

### 補足
- Go 1.21: `ErrJSTemplate` 発生
- Go 1.22+: 同じ構文で成功（回避コードあり）

---

## 🔹 UC6. ファイルアップロードとセーフFSアクセス

### 概要
JSONなどのファイルをアップロードし、`os.Root` を用いてファイルパスのセーフアクセスを試行。

### API仕様
- `POST /upload`: ファイルアップロード（multipart/form-data）
- `GET /uploads/{filename}`: 安全なファイル参照（os.Root）

### ユースケース構造
- `UploadFileUsecase`, `GetUploadedFileUsecase`
- `FSBox`（ドメインサービス or インフラ抽象）

### 補足
- `os.Root.Open(filename)` で限定範囲にアクセス制限
- `ServeFileFS` による埋め込みFS対応あり

---

## 🔹 UC7. HTTP/2・TLS・H2C設定の確認

### 概要
通信プロトコルとTLSの状態を表示し、暗号方式、クライアント情報、プロトコルネゴシエーション状況を確認。

### API仕様
- `GET /tls/info`: TLS・HTTP2ステータス表示
- `GET /tls/raw`: 生クライアントヘッダ＆TLSセッション表示

### ユースケース構造
- `TLSInfoResolver`（ドメインサービス）
- `TLSState`（VO）

### 補足
- `httptrace` による詳細取得
- `Server.Protocols`, `Transport.Protocols` の差異確認
- ECH・X25519・Kyber768等の対応表示（Go 1.23+）

---

## 🔹 UC8. PGO（Profile-Guided Optimization）テスト

### 概要
アプリケーションにプロファイルを適用してビルドし、性能がどう変わるかを比較・可視化。

### ページ
- `GET /pgo/stats`: PGO前後の比較実行時間を表示

### 内容
- `go build -pgo=auto`
- `testing.Benchmark` による処理時間の比較表示
- 利用プロファイルのサンプル提示

---

各ユースケースは Clean Architecture に沿って、以下の対応関係を持ちます：

- **Controller**: HTTPルーティング・リクエストバインディング
- **Usecase (InputPort)**: 要件ロジック（リクエスト→レスポンス変換）
- **OutputPort (Gateway)**: 永続化や外部との接続抽象
- **Domain Model/Service**: 不変な業務ルール・ロジック

---
