# loglinter

Go linter for log message validation. Checks log calls from `log/slog` and `go.uber.org/zap` against a set of rules.

## Rules

| Rule | Description | Example violation |
|------|-------------|-------------------|
| `lowercase` | Message must start with a lowercase letter | `slog.Info("Starting server")` |
| `english` | Message must be in English only | `slog.Info("запуск сервера")` |
| `special_chars` | No special characters or emoji | `slog.Info("server started!🚀")` |
| `sensitive_keywords` | No sensitive data in arguments | `slog.Info("auth", password)` |
| `custom_patterns` | User-defined regex patterns | AWS keys, JWT tokens, etc. |

## Project structure

```
loglinter/
├── cmd/
│   ├── linter/
│   │   └── main.go          # standalone binary
│   └── plugin/
│       └── main.go          # golangci-lint plugin
├── config/
│   ├── config.go
│   └── config.yaml          # example config
├── pkg/
│   └── analyzer/
│       ├── analyzer.go      # analysis.Analyzer definition
│       ├── checker.go       # AST traversal, logger detection
│       ├── rules.go         # rule implementations
│       ├── rules_test.go
│       ├── analyzer_test.go
│       ├── testdata/        # analysistest data
│       └── detector/
│           ├── detector.go  # regex pattern detector
│           └── detector_test.go
└── .golangci.yml
```

## Requirements

- Go 1.26+
- `CGO_ENABLED=1` (required for the golangci-lint plugin)
- golangci-lint v2.x

## Build

### Standalone binary

```bash
go build -o loglinter ./cmd/linter/
```

Run against a package:

```bash
CONFIG_PATH=./config/config.yaml ./selectellinter ./...
```

### golangci-lint plugin

The plugin and golangci-lint **must be built with the same Go version and dependencies**.

> **Note:** The `.so` plugin requires Linux and `CGO_ENABLED=1`. The official golangci-lint binary is built without CGO, so the plugin must be used with a locally installed golangci-lint.

```bash
# Check which Go version golangci-lint was built with
go version -m $(which golangci-lint)

# Build the plugin
CGO_ENABLED=1 go build -buildmode=plugin -o loglint.so ./cmd/plugin/
```

## Configuration

Set the config path via `CONFIG_PATH` environment variable:

```bash
CONFIG_PATH=./config/config.yaml golangci-lint run ./...
```

`config/config.yaml` example:

```yaml
rules:
  lowercase: true
  english: true
  special_chars: true
  sensitive_keywords: true
  custom_patterns: true

auto_fix:
  enabled: true
  lowercase: true
  special_chars: true
  sensitive_keywords: true
  custom_patterns: true

sensitive_keywords:
  - password
  - token
  - api_key
  - secret
  - credit_card

custom_patterns:
  - name: aws_key
    pattern: "(?i)AKIA[0-9A-Z]{16}"
    message: "AWS access key detected"
    auto_fix: true
    replacement: "AWS_KEY /REDACTED/"

  - name: jwt_token
    pattern: "eyJ[a-zA-Z0-9_-]+\\.[a-zA-Z0-9_-]+\\.[a-zA-Z0-9_-]+"
    message: "JWT token detected"
    auto_fix: true
    replacement: "JWT /REDACTED/"

  - name: password_inline
    pattern: "(?i)password\\s*[:=]\\s*\\S+"
    message: "inline password value detected"
    auto_fix: true
    replacement: "password /REDACTED/"
```

## golangci-lint integration

`.golangci.yml`:

```yaml
version: "2"

linters:
  enable:
    - loglint
  settings:
    custom:
      loglint:
        path: ./loglint.so
        description: "Log message linter"
        original-url: selectellinter
```

Run:

```bash
CONFIG_PATH=./config/config.yaml golangci-lint run ./...
```

## Usage example

Given `examples/ex.go`:

```go
package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

func main() {
	// --- slog logger via variable ---
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting server on port 8080")       // ❌ uppercase
	logger.Error("ошибка подключения к базе данных")  // ❌ not English
	logger.Warn("server started!!! 🚀")               // ❌ special chars + emoji
	password := "supersecret"
	logger.Info("user login", "pass", password)       // ❌ sensitive keyword in ident
	logger.Info("password: supersecret")              // ❌ inline password pattern

	// --- zap logger via variable ---
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	zapLogger.Info("Failed to process request")       // ❌ uppercase
	zapLogger.Error("запуск воркера завершён")         // ❌ not English
	zapLogger.Warn("connection lost... 💀")           // ❌ special chars + emoji
	userToken := "eyJhbGciOiJIUzI1NiJ9.payload.sig"
	zapLogger.Info("request received", zap.String("auth", userToken)) // ❌ sensitive keyword in ident
}
```

Running the linter:

```
$ CONFIG_PATH=./config/config.yaml golangci-lint run ./examples/ex.go

ex.go:14:14: linter: Log must start with a lowercase letter (loglint)
        logger.Info("Starting server on port 8080")
                    ^
ex.go:16:2:  linter: log message must be in English (found: о) (loglint)
        logger.Error("ошибка подключения к базе данных")
        ^
ex.go:18:14: linter: log message must not contain special characters or emoji (loglint)
        logger.Warn("server started!!! 🚀")
                    ^
ex.go:21:2:  linter: log message contains sensitive keyword: "password" (loglint)
        logger.Info("user login", "pass", password)
        ^
ex.go:23:14: linter: inline password value detected (loglint)
        logger.Info("password: supersecret")
                    ^
ex.go:28:17: linter: Log must start with a lowercase letter (loglint)
        zapLogger.Info("Failed to process request")
                       ^
ex.go:30:2:  linter: log message must be in English (found: з) (loglint)
        zapLogger.Error("запуск воркера завершён")
        ^
ex.go:32:17: linter: log message must not contain special characters or emoji (loglint)
        zapLogger.Warn("connection lost... 💀")
                       ^
ex.go:35:2:  linter: log message contains sensitive keyword: "token" (loglint)
        zapLogger.Info("request received", zap.String("auth", userToken))
        ^

9 issues: loglint: 9
```

Running with `--fix` applies all auto-fixable suggestions:

```
$ CONFIG_PATH=./config/config.yaml golangci-lint run --fix ./examples/ex.go
```

After auto-fix the file becomes:

```go
logger.Info("starting server on port 8080")  // ✅ first letter lowercased
logger.Warn("server started ")               // ✅ special chars and emoji removed
logger.Info("password /REDACTED/")           // ✅ inline password redacted by pattern

zapLogger.Info("failed to process request")  // ✅ first letter lowercased
zapLogger.Warn("connection lost ")           // ✅ special chars and emoji removed
```

Remaining issues after `--fix` require **manual intervention**:

| Issue | Why auto-fix is not possible |
|-------|------------------------------|
| `log message must be in English` | Cannot translate text statically |
| `log message contains sensitive keyword` (via variable name) | Cannot rename variables or remove arguments |

## Run tests

```bash
go test ./...
```

## Platform support

The golangci-lint `.so` plugin works **Linux only** (Go limitation for `buildmode=plugin`).

For macOS and Windows use the standalone binary via `go vet`:

```bash
go build -o loglinter ./cmd/linter/
CONFIG_PATH=./config/config.yaml go vet -vettool=./loglinter ./...
```
