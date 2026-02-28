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
│   ├── main/
│   │   └── main.go          # standalone binary
│   └── plugin/
│       └── plugin.go        # golangci-lint plugin
├── config/
│   ├── config.go
│   └── config.yaml          # example config
├── pkg/
│   └── analyzer/
│       ├── analyzer.go      # analysis.Analyzer definition
│       ├── checker.go       # AST traversal, logger detection
│       ├── rules.go         # rule implementations
│       └── detector/
│           └── detector.go  # regex pattern detector
├── testdata/
│   ├── lowercase/
│   ├── english/
│   ├── specialchars/
│   ├── sensitive/
│   └── custom/          
└── .golangci.yml
```

## Requirements

- Go 1.26.0
- `CGO_ENABLED=1` (required for the golangci-lint plugin)
- golangci-lint v2.x

## Build

### Standalone binary

```bash
go build -buildmode=plugin -o loglint.so ./cmd/plugin/
```

Run against a package:

```bash
CONFIG_PATH=./config.yaml ./selectellinter ./...
```

### golangci-lint plugin

The plugin and golangci-lint **must be built with the same Go version and dependencies**.

```bash
# Check which Go version golangci-lint was built with
go version -m $(which golangci-lint)

# Build the plugin
CGO_ENABLED=1 go build -buildmode=plugin -o loglint.so ./cmd/plugin/
```

## Configuration

Create a `local.yaml` (or any path, set via `CONFIG_PATH` env variable):

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
  english: true
  special_chars: true
  sensitive_keywords: true
  custom_patterns: true

sensitive_keywords:
  - password
  - token
  - apiKey
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

Add to `.golangci.yml`:

```yaml
version: "2"

linters:
  settings:
    custom:
      loglint:
        path: ./loglint.so
        description: "Log message linter"
        original-url: selectellinter
        settings:
          config_path: ./config.yaml

linters:
  enable:
    - loglint
```

Set the config path before running:

```bash
CONFIG_PATH=./config.yaml golangci-lint run ./...
```

## Usage examples

```go
// ❌ will be reported
slog.Info("Starting server on port 8080")      // Log must start with a lowercase letter
slog.Info("запуск сервера")                     // log message must be in English
slog.Error("connection failed!!!")              // log message must not contain special characters or emoji
slog.Debug("auth", "api_key", apiKey)           // log message contains sensitive keyword: "apiKey"
slog.Info("key: AKIAIOSFODNN7EXAMPLE")          // AWS access key detected

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

When `auto_fix.enabled: true`, golangci-lint will suggest fixes that can be applied with:

```bash
CONFIG_PATH=./config.yaml golangci-lint run --fix ./...
```

Example:
```go
// before
slog.Info("Starting server")

// after auto-fix
slog.Info("starting server")
```