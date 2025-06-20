# Go Web Console — プロジェクト計画書（Part 1: 概要・背景・システム構成）

---

## 1. 背景・目的

Go 1.21〜1.24 にかけて、Webアプリケーション開発に大きな影響を与える言語機能および標準ライブラリの拡張が数多く追加されました。

特に注目すべきは以下の点です：

- `net/http.ServeMux` のルーティング強化（パターン優先順位、HTTPメソッド、パス変数）
- `log/slog` による構造化ログの標準化
- `context` のキャンセル・締切理由の明確化（`Cause`, `WithoutCancel`, `AfterFunc`）
- `html/template` の表現力強化（整数 `range`, JSリテラル埋め込み対応）
- HTTP/2 / H2C / TLS構成の制御・安全性向上
- PGO（Profile-Guided Optimization）によるパフォーマンス最適化
- `os.Root`, `ServeFileFS`, `unique`, `AddCleanup`, `weak` などの実用系追加機能

これらを体系的に試し、ベストプラクティスを学習・伝播可能な形で提供する**自己ホスト型の開発者向けコンソールアプリ**が求められています。

---

## 2. アプリケーション概要

### 📛 名称
**Go Web Console**

### 🧑‍💻 利用者像
- Go Web/API 開発者（初中級〜上級）
- クリーンアーキテクチャに基づくアプリ設計を学びたい人
- Goの最新機能を使って性能やセキュリティを向上させたい人

### 🎯 目的
- Web開発における **Goの最新機能を体験・理解・評価** する
- DDD + Clean Architecture に沿って **堅牢かつ拡張可能な設計を示す**
- 開発／検証／教育用途に使える「**自己診断可能なWebサーバ**」の提供

---

## 3. システム構成方針（Clean Architecture + DDD）

### 🏗️ アーキテクチャ階層

```
Clean Architecture の4層構造 + DDDによるドメイン分離
----------------------------------------------------
[外部層]        → WebUI / CLI / Tests 等
[インターフェース層] Controller / Gateway / View
[ユースケース層] Usecases / Interactor
[ドメイン層]     Domain Model / Entity / Service / ValueObject
[インフラ層]     DB, ファイル, ネットワーク, FS, TLS, FSBox 等
```

---

## 4. ディレクトリ構成（標準構成例）

```
go-web-console/
├── cmd/                    # エントリポイント
│   └── server/            # main.go / DI / 設定読み込み
├── internal/
│   ├── domain/            # ドメイン層
│   │   ├── model/         # エンティティ・VO
│   │   └── service/       # ドメインサービス
│   ├── usecase/           # ユースケース層（InputPort / Interactor）
│   ├── interface/         # インターフェース層
│   │   ├── controller/    # HTTPハンドラ群
│   │   └── gateway/       # ファイル/ログ/DBの抽象化（OutputPort）
│   ├── infrastructure/    # 実装技術依存（os, slog, fs 等）
│   └── shared/            # 共通ユーティリティ（context/logエイリアス等）
├── templates/             # HTMLテンプレート（embed対象）
├── static/                # 静的ファイル（JS/CSS/アイコン等）
├── config/                # 設定ファイル・pgo.prof
├── pgo.prof               # PGO用プロファイルファイル
├── go.mod
└── main.go
```

---

## 5. モジュールの責務分離（例）

| ディレクトリ | 主な役割 |
|--------------|----------|
| `domain/model` | `LogEntry`, `User`, `RequestTrace` などの純粋なビジネスモデル |
| `domain/service` | `CookieParserService`, `TLSInfoResolver` などのドメインサービス |
| `usecase/` | `GetLogsUsecase`, `RestartTaskUsecase` など、ユーザー視点の操作単位の業務処理 |
| `interface/controller` | `LogHandler`, `TemplateTestHandler` など HTTP レベルの入出力定義 |
| `interface/gateway` | `FileLogStore`, `CookieReader`, `HeaderTracer` などの外部IO抽象 |
| `infrastructure/` | `log/slog`, `fs.Embed`, `os.Root` などの具体実装 |
| `shared/` | `ctxutil`, `errorutil`, `httperr`, `slogx` 等の共通処理 |

---

## 6. 実行時構成と依存構成

- **サーバ起動：** `cmd/server/main.go` から DI コンテナ（go.uber.org/fx または独自手動注入）で全ての依存解決
- **テンプレート・静的配信：** `embed.FS` を使って `templates/`, `static/` をバイナリ埋込
- **設定：** YAML or JSON による `config/config.yaml` 読み込み（PGO, TLS, ServeMuxの動作トグルなど含む）
- **ログ：** `slog` ハンドラを JSON / text 出力に切替可能

---
