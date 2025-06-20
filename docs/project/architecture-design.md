# Go Web Console — プロジェクト計画書（Part 3: ドメイン/ユースケース/インターフェース設計）

---

## 1. ドメインモデル（Entity / VO）

### 📄 `domain/model/log.go`
```go
package model

import "time"

type LogLevel string

const (
    InfoLevel  LogLevel = "info"
    WarnLevel  LogLevel = "warn"
    ErrorLevel LogLevel = "error"
)

type LogEntry struct {
    Timestamp time.Time
    Level     LogLevel
    Message   string
    Fields    map[string]any
}
```

---

### 📄 `domain/model/context_trace.go`
```go
package model

import "context"

type ContextTrace struct {
    ID        string
    StartTime time.Time
    Deadline  time.Time
    Cancelled bool
    Cause     error
}

func NewContextTrace(ctx context.Context) *ContextTrace {
    dl, _ := ctx.Deadline()
    return &ContextTrace{
        ID:        "generated-request-id", // TODO: UUID
        StartTime: time.Now(),
        Deadline:  dl,
    }
}
```

---

### 📄 `domain/service/tls_info_resolver.go`
```go
package service

import (
    "crypto/tls"
)

type TLSInfo struct {
    Version       string
    CipherSuite   string
    NegotiatedALPN string
    ECHAccepted   bool
}

func ResolveTLSInfo(state *tls.ConnectionState) TLSInfo {
    return TLSInfo{
        Version:       tlsVersionName(state.Version),
        CipherSuite:   tls.CipherSuiteName(state.CipherSuite),
        NegotiatedALPN: state.NegotiatedProtocol,
        ECHAccepted:   state.ECHAccepted,
    }
}
```

---

## 2. ユースケース（InputPort / Interactor）

### 📄 `usecase/log/viewer.go`
```go
package log

import (
    "context"
    "go-web-console/internal/domain/model"
)

type ViewerUsecase interface {
    GetLogs(ctx context.Context, level model.LogLevel) ([]model.LogEntry, error)
}
```

### 📄 `usecase/log/viewer_interactor.go`
```go
package log

import (
    "context"
    "go-web-console/internal/domain/model"
    "go-web-console/internal/interface/gateway"
)

type viewerInteractor struct {
    repo gateway.LogRepository
}

func NewViewerUsecase(repo gateway.LogRepository) ViewerUsecase {
    return &viewerInteractor{repo}
}

func (u *viewerInteractor) GetLogs(ctx context.Context, level model.LogLevel) ([]model.LogEntry, error) {
    return u.repo.FindByLevel(ctx, level)
}
```

---

## 3. インターフェース層（Controller / Gateway）

### 📄 `interface/controller/log_handler.go`
```go
package controller

import (
    "encoding/json"
    "net/http"
    "go-web-console/internal/domain/model"
    "go-web-console/internal/usecase/log"
)

type LogHandler struct {
    uc log.ViewerUsecase
}

func NewLogHandler(uc log.ViewerUsecase) *LogHandler {
    return &LogHandler{uc}
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    level := r.PathValue("level")
    logs, err := h.uc.GetLogs(r.Context(), model.LogLevel(level))
    if err != nil {
        http.Error(w, "failed to get logs", http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(logs)
}
```

---

### 📄 `interface/gateway/log_repository.go`
```go
package gateway

import (
    "context"
    "go-web-console/internal/domain/model"
)

type LogRepository interface {
    FindByLevel(ctx context.Context, level model.LogLevel) ([]model.LogEntry, error)
}
```

### 📄 `infrastructure/log/file_logger.go`
```go
package loginfra

import (
    "context"
    "go-web-console/internal/domain/model"
)

type FileLogStore struct {
    filePath string
}

func NewFileLogStore(path string) *FileLogStore {
    return &FileLogStore{path}
}

func (f *FileLogStore) FindByLevel(ctx context.Context, level model.LogLevel) ([]model.LogEntry, error) {
    // JSONファイルを読み込んでレベルでフィルタ
    return nil, nil // 実装は省略
}
```

---

## 4. その他ユーティリティ（Shared）

### 📄 `shared/context/ctxutil.go`
```go
package ctxutil

import "context"

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, ctxKey("requestID"), id)
}

func GetRequestID(ctx context.Context) string {
    v := ctx.Value(ctxKey("requestID"))
    if s, ok := v.(string); ok {
        return s
    }
    return ""
}

type ctxKey string
```

---

## 補足

- DIは `cmd/server/main.go` で `controller -> usecase -> gateway` を手動で注入（Wire/fx利用も可能）
- ドメイン層は外部依存を持たない純粋な型定義
- インフラの組み換え（e.g., file → stdout, embed.FS）も可能

---
