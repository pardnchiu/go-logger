# Go Logger (Golang)

> A Logging package for Golang with automatic rotation, multi-level log classification, and file management, featuring comprehensive error handling mechanisms.<br>
> Primarily designed for use in `pardnchiu/go-*` packages

[![license](https://img.shields.io/github/license/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) 
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/releases) 
[![readme](https://img.shields.io/badge/readme-中文-blue)](https://github.com/pardnchiu/go-logger/blob/main/README.zh.md) 

## Three Key Features

- **Tree Structure & slog Standardized Output Support**: JSON format compatible with Go's standard log/slog package for structured logging, Text format uses tree structure to enhance readable experience
- **Complete Multi-level Log Classification**: Supports 8 levels (DEBUG, TRACE, INFO, NOTICE, WARNING, ERROR, FATAL, CRITICAL)
- **Automatic File Rotation & Cleanup**: Automatically rotates and creates backups when file size reaches limit, intelligently cleans expired files to maintain configured backup count

## How to use

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
  
  goLogger "github.com/pardnchiu/go-logger"
)

func main() {
  // Create configuration
  config := &goLogger.Log{
    Path:      "./logs",              // Log directory
    Stdout:    true,                  // Output to terminal as well
    MaxSize:   16 * 1024 * 1024,      // 16MB file size limit
    MaxBackup: 5,                     // Keep 5 backup files
    Type:      "json",                // "json" for slog standard, "text" for tree format
  }
  
  // Initialize logger
  logger, err := goLogger.New(config)
  if err != nil {
    panic(err)
  }
  defer logger.Close()
  
  // Use different levels to log messages
  logger.Debug("This is debug message", "Detailed debug information")
  logger.Trace("Tracing program execution flow")
  logger.Info("General information message")
  logger.Notice("Message that needs attention")
  logger.Warn("Warning message")
  
  // Error handling
  err = errors.New("An error occurred")
  logger.WarnError(err, "Warning message")
  logger.Error(err, "Additional message when handling error")
  logger.Fatal(err, "Critical error")
  logger.Critical(err, "System critical error")
  
  // Manual flush cache
  logger.Flush()
}
```

### Configuration Details

```go
type Log struct {
  Path      string // Log file directory path (default: ./logs)
  Stdout    bool   // Whether to output to standard output (default: false)
  MaxSize   int64  // Maximum log file size in bytes (default: 16MB)
  MaxBackup int    // Maximum number of backup files (default: 5)
  Type      string // Output format: "json" for slog standard, "text" for tree format (default: "text")
}
```

## Output Formats

### JSON Format (slog Standard)
When `Type: "json"`, logs are output in slog-compatible structured format:

```json
{"timestamp":"2024/01/15 14:30:25.123456","level":"INFO","message":"Application started","data":null}
{"timestamp":"2024/01/15 14:30:25.123457","level":"ERROR","message":"Database connection failed","data":["Connection timeout","Retrying in 5 seconds"]}
```

Benefits:
- Compatible with Go's standard `log/slog` package
- Machine-readable structured logging
- Easy integration with log aggregation tools
- Consistent JSON schema across all log levels

### Text Format (Tree Structure)
When `Type: "text"`, logs are displayed in human-readable tree format:

```
2024/01/15 14:30:25.123456 Application started
2024/01/15 14:30:25.123457 [ERROR] Database connection failed
2024/01/15 14:30:25.123457 ├── Connection timeout
2024/01/15 14:30:25.123457 └── Retrying in 5 seconds
```

Benefits:
- Human-friendly visual representation
- Clear hierarchical message structure
- Enhanced readability for debugging

## Log Level Description

### Debug and Trace
Logged to `debug.log`
```go
logger.Debug("Variable values", "x = 10", "y = 20")
logger.Trace("Function call", "Starting user request processing")
```

### Info, Notice, Warning
Logged to `output.log`
```go
logger.Info("Application started")           // No prefix
logger.Notice("Configuration file reloaded") // [NOTICE] prefix
logger.Warn("Memory usage too high")         // [WARNING] prefix
```

### Error, Fatal, Critical
Logged to `error.log`
```go
logger.Error(err, "Retry attempt 3")        // [ERROR] prefix
logger.Fatal(err, "Cannot start service")   // [FATAL] prefix
logger.Critical(err, "System crash")        // [CRITICAL] prefix
```

## Core Features

### Logger Management

- **New** - Create new logger instance
  ```go
  logger, err := goLogger.New(config)
  ```
  - Initialize log directory, ensure path exists
  - Open three log files: debug.log, output.log, error.log
  - Configure log handlers for each level
  - Check existing file sizes, perform rotation if necessary

- **Close** - Safely close logger
  ```go
  err := logger.Close()
  ```
  - Close all open file handles
  - Mark logger as closed
  - Ensure no resource leaks

- **Flush** - Force write cache
  ```go
  err := logger.Flush()
  ```
  - Write all cached log content to disk
  - Ensure important logs are not lost

### File Rotation Mechanism

#### Automatic Rotation
- Check file size before each log write
- Automatically rotate when exceeding `MaxSize` limit
- Backup file naming format: `filename.YYYYMMDD_HHMMSS`

#### Backup Management
- Keep the latest `MaxBackup` backup files
- Automatically delete expired old backups
- Sort by modification time, keep the newest files

### Concurrency Safety Mechanism

#### Read-Write Lock Protection
- Use `sync.RWMutex` to protect critical sections
- Write operations acquire write lock, ensuring atomicity
- Read operations use read lock, improving concurrent performance

## Usage Examples

### Basic Logging with JSON Format
```go
logger, _ := goLogger.New(&goLogger.Log{
  Path:    "./logs",
  Stdout:  true,
  MaxSize: 1024 * 1024, // 1MB
  Type:    "json",      // slog standard format
})
defer logger.Close()

logger.Info("Application started")
logger.Debug("Loading configuration file", "config.json")
logger.Warn("Memory usage", "85%")
```

### Basic Logging with Text Format
```go
logger, _ := goLogger.New(&goLogger.Log{
  Path:    "./logs",
  Stdout:  true,
  MaxSize: 1024 * 1024, // 1MB
  Type:    "text",      // Tree structure format
})
defer logger.Close()

logger.Info("Application started")
logger.Debug("Loading configuration file", "config.json")
logger.Warn("Memory usage", "85%")
```

### Error Handling
```go
if err := connectDatabase(); err != nil {
  logger.Error(err, "Database connection failed", "Retrying...")
  return logger.Fatal(err, "Cannot establish database connection")
}
```

### Concurrent Environment
```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
    defer wg.Done()
    logger.Info(fmt.Sprintf("Goroutine %d running", id))
  }(i)
}

wg.Wait()
logger.Flush() // Ensure all logs are written
```

## License

This source code project is licensed under the [MIT](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) License.

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