> [!Note]
> This content is translated by LLM. Original text can be found [here](README.zh.md)

# Go Logger

> A Golang logging package with automatic rotation, multi-level log classification, file management capabilities, and comprehensive error handling mechanisms.<br>
> Primarily designed for use in `pardnchiu/go-*` packages

[![lang](https://img.shields.io/badge/lang-Go-blue)](README.zh.md) 
[![license](https://img.shields.io/github/license/pardnchiu/go-logger)](LICENSE)
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/releases)
![card](https://goreportcard.com/badge/github.com/pardnchiu/go-logger)<br>
[![readme](https://img.shields.io/badge/readme-EN-white)](README.md)
[![readme](https://img.shields.io/badge/readme-ZH-white)](README.zh.md) 

## Three Core Features

### Support for slog Standardization and Tree Structure Output
JSON uses Go's standard `log/slog` package for structured logging
Text adopts tree structure to enhance readability

### Complete Multi-Level Log Classification
Supports 8 levels (`DEBUG`, `TRACE`, `INFO`, `NOTICE`, `WARNING`, `ERROR`, `FATAL`, `CRITICAL`)

### Automatic File Rotation and Cleanup
Automatically rotates and creates backups when files reach size limits, intelligently cleans expired files to maintain configured backup count

## Usage

### Installation
```bash
go get github.com/pardnchiu/go-logger
```

### Initialization
```go
package main

import (
  "fmt"
  "errors"
  
  "github.com/pardnchiu/go-logger"
)

func main() {
  config := &goLogger.Log{
    Path:      "./logs",              // Log directory
    Stdout:    true,                  // Also output to terminal
    MaxSize:   16 * 1024 * 1024,      // 16MB file size limit
    MaxBackup: 5,                     // Keep 5 backup files
    Type:      "json",                // "json" for slog standard, "text" for tree format
  }
  
  // Initialize
  logger, err := goLogger.New(config)
  if err != nil {
    panic(err)
  }
  defer logger.Close()
  
  // Log messages at different levels
  logger.Debug("This is debug message", "detailed debug info")
  logger.Trace("Trace program execution flow")
  logger.Info("General information message")
  logger.Notice("Message that needs attention")
  logger.Warn("Warning message")
  
  // Error handling
  err = errors.New("an error occurred")
  logger.WarnError(err, "Warning message for handling error")
  logger.Error(err, "Additional message for handling error")
  logger.Fatal(err, "Critical error")
  logger.Critical(err, "System critical error")
  
  // Flush cache
  logger.Flush()
}
```

## Configuration

```go
type Log struct {
  Path      string // Log file directory path (default: ./logs)
  Stdout    bool   // Whether to output to stdout (default: false)
  MaxSize   int64  // Maximum log file size in bytes (default: 16MB)
  MaxBackup int    // Maximum number of backup files (default: 5)
  Type      string // Output format: "json" for slog standard, "text" for tree format (default: "text")
}
```

## Output Formats

### slog Standard
When `Type: "json"`, logs are output in `log/slog` structured format:

```json
{"timestamp":"2024/01/15 14:30:25.123456","level":"INFO","message":"Application started","data":null}
{"timestamp":"2024/01/15 14:30:25.123457","level":"ERROR","message":"Database connection failed","data":["Connection timeout","Retry in 5 seconds"]}
```
- Directly uses Go's standard `log/slog` package
- Easy integration with log aggregation tools
- Consistent JSON schema across all log levels

### Tree Structure
When `Type: "text"`, logs are displayed in tree format:

```
2024/01/15 14:30:25.123457 [ERROR] Database connection failed
2024/01/15 14:30:25.123457 ├── Connection timeout
2024/01/15 14:30:25.123457 └── Retry in 5 seconds
```
- Clear hierarchical message structure
- Enhanced readability during debugging

## Log Levels

### Debug and Trace
Logged to `debug.log`
```go
logger.Debug("Variable values", "x = 10", "y = 20")
logger.Trace("Function call", "Started processing user request")
```

### Info, Notice, Warning
Logged to `output.log`
```go
logger.Info("Application started")                    // No prefix
logger.Notice("Configuration file reloaded")          // [NOTICE] prefix
logger.Warn("Memory usage is high")                   // [WARNING] prefix
logger.WarnError(err, "Non-system-affecting error")   // [WARNING] prefix
```

### Error, Fatal, Critical
Logged to `error.log`
```go
logger.Error(err, "Retry attempt 3")         // [ERROR] prefix
logger.Fatal(err, "Unable to start service") // [FATAL] prefix
logger.Critical(err, "System crash")         // [CRITICAL] prefix
```

## Available Functions

- **New** - Create a new logger instance
  ```go
  logger, err := goLogger.New(config)
  ```
  - Initialize log directory, ensure path exists
  - Initialize three log files: `debug.log`, `output.log`, `error.log`
  - Set up log handlers for each level

- **Close** - Properly close the logger
  ```go
  err := logger.Close()
  ```
  - Mark logger as closed
  - Ensure no resource leaks

- **Flush** - Force write to files
  ```go
  err := logger.Flush()
  ```
  - Write all cached log content to disk
  - Ensure logs are not lost

### File Rotation Mechanism

#### Automatic Rotation
- Check file size before each log write
- Automatically rotate when exceeding `MaxSize` limit
- Backup file naming format: `filename.YYYYMMDD_HHMMSS`

#### Backup Management
- Keep the latest `MaxBackup` backup files
- Automatically delete expired old backups
- Sort by modification time, keep the newest files

## License

This project is licensed under the [MIT](LICENSE) License.

## Author

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
  <img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
  <img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

***

©️ 2025 [邱敬幃 Pardn Chiu](https://pardn.io)