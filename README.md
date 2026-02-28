# selectellinter

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
selectellinter/
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

- Go 1.22+
- `CGO_ENABLED=1` (required for the golangci-lint plugin)
- golangci-lint v2.x

## Build

### Standalone binary

```bash
go build -o selectellinter ./cmd/linter/
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
    replacement: "AWS_KEY=[REDACTED]"

  - name: jwt_token
    pattern: "eyJ[a-zA-Z0-9_-]+\\.[a-zA-Z0-9_-]+\\.[a-zA-Z0-9_-]+"
    message: "JWT token detected"
    auto_fix: true
    replacement: "JWT=[REDACTED]"
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

## Usage examples

```go
// ❌ will be reported
slog.Info("Starting server on port 8080")   // Log must start with a lowercase letter
slog.Info("запуск сервера")                  // log message must be in English
slog.Error("connection failed!!!")           // log message must not contain special characters or emoji
slog.Debug("request", "api_key", apiKey)     // log message contains sensitive keyword: "api_key"
slog.Info("key: AKIAIOSFODNN7EXAMPLE")       // AWS access key detected

// ✅ correct
slog.Info("starting server on port 8080")
slog.Info("starting server")
slog.Error("connection failed")
slog.Debug("request completed")
```

## Run tests

```bash
go test ./...
```

## Auto-fix

When `auto_fix.enabled: true`, suggested fixes can be applied with:

```bash
CONFIG_PATH=./config/config.yaml golangci-lint run --fix ./...
```

```go
// before
slog.Info("Starting server")

// after auto-fix
slog.Info("starting server")
```

## Platform support

The golangci-lint `.so` plugin works **Linux only** (Go limitation for `buildmode=plugin`).

For macOS and Windows use the standalone binary via `go vet`:

```bash
go build -o selectellinter ./cmd/linter/
CONFIG_PATH=./config/config.yaml go vet -vettool=./selectellinter ./...
```
