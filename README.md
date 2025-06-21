# Go Web Console

Go Web Console is a sample project showcasing modern Go web development
practices. It demonstrates a clean architecture approach with manual
dependency injection and uses Go 1.24 features such as the improved
`net/http` `ServeMux`, context tracing, and the `slog` structured logging
package.

## Requirements

- Go 1.24 or later

## Setup

1. Install dependencies (internet access required):
   ```bash
   go mod tidy
   ```
2. Run the HTTP server:
   ```bash
   go run ./cmd/server
   ```

## Directory Overview

```
cmd/                     Entry points
  server/                main.go and DI setup
internal/                Application packages
  domain/                Entities and domain services
    model/
    service/
  usecase/               Business use cases
  interface/             HTTP controllers and gateways
    controller/
    gateway/
  infrastructure/        Technical implementations (fs, slog, etc.)
  shared/                Cross-cutting utilities
templates/               HTML templates
static/                  Static files (JS/CSS)
config/                  Configuration files
```

For more details see the [project documents](docs/project/).
