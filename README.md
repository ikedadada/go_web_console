# Go Web Console

Go Web Console is an example web server demonstrating new Go features (1.21+)
and clean architecture practices. It targets developers who want to evaluate the
latest language capabilities while learning about layered application design.

## Requirements

- Go 1.22 or later

## Setup

```bash
# Clone and enter the repository
git clone <repo-url>
cd go_web_console

# Download dependencies
go mod download
```

## Usage

Build and run the server:

```bash
go run ./cmd/server
```

Run the full test suite:

```bash
go test ./...
```

## Directory Layout

```
go-web-console/
├── cmd/                    # entry points
│   └── server/            # main.go & DI setup
├── internal/
│   ├── domain/            # business entities & services
│   ├── usecase/           # application use cases
│   ├── interface/         # HTTP controllers & gateways
│   ├── infrastructure/    # concrete implementations
│   └── shared/            # common utilities
├── templates/             # HTML templates
├── static/                # JS/CSS assets
├── config/                # configuration files
├── pgo.prof               # profile for PGO
├── go.mod
└── main.go
```

This repository currently contains documentation outlining goals and design.
Future commits will add the implementation.
