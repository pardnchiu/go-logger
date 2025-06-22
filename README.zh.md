# Go Logger (Golang)

> 一個 Golang 日誌套件，具備自動輪替、多層級日誌分類和檔案管理功能，以及完整的錯誤處理機制。<br>
> 主要設計用於 `pardnchiu/go-*` 套件中

[![license](https://img.shields.io/github/license/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) 
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/releases) 
[![readme](https://img.shields.io/badge/readme-中文-blue)](https://github.com/pardnchiu/go-logger/blob/main/README.zh.md) 

## 三大主軸

- **支持樹狀結構與 slog 標準化輸出**：JSON 相容 Go 標準 log/slog 套件用於結構化記錄，Text 採用樹狀結構提升閱讀體驗
- **完整多層級日誌分類**：支援 8 個層級（DEBUG、TRACE、INFO、NOTICE、WARNING、ERROR、FATAL、CRITICAL）
- **自動檔案輪替與清理**：檔案達大小限制時自動輪替並建立備份，智慧清理過期檔案維護設定的備份數量

## 使用方法

### 安裝
```bash
go get github.com/pardnchiu/go-logger
```

### 初始化
```go
package main

import (
  "fmt"
  "errors"
  
  goLogger "github.com/pardnchiu/go-logger"
)

func main() {
  // 建立設定
  config := &goLogger.Log{
    Path:      "./logs",              // 日誌目錄
    Stdout:    true,                  // 同時輸出到終端
    MaxSize:   16 * 1024 * 1024,      // 16MB 檔案大小限制
    MaxBackup: 5,                     // 保留 5 個備份檔案
    Type:      "json",                // "json" 為 slog 標準，"text" 為樹狀格式
  }
  
  // 初始化日誌記錄器
  logger, err := goLogger.New(config)
  if err != nil {
    panic(err)
  }
  defer logger.Close()
  
  // 使用不同層級記錄訊息
  logger.Debug("這是除錯訊息", "詳細除錯資訊")
  logger.Trace("追蹤程式執行流程")
  logger.Info("一般資訊訊息")
  logger.Notice("需要注意的訊息")
  logger.Warn("警告訊息")
  
  // 錯誤處理
  err = errors.New("發生錯誤")
  logger.Error(err, "處理錯誤時的附加訊息")
  logger.Fatal(err, "嚴重錯誤")
  logger.Critical(err, "系統關鍵錯誤")
  
  // 手動清空快取
  logger.Flush()
}
```

### 設定詳細說明

```go
type Log struct {
  Path      string // 日誌檔案目錄路徑（預設：./logs）
  Stdout    bool   // 是否輸出到標準輸出（預設：false）
  MaxSize   int64  // 日誌檔案最大大小（位元組）（預設：16MB）
  MaxBackup int    // 最大備份檔案數量（預設：5）
  Type      string // 輸出格式："json" 為 slog 標準，"text" 為樹狀格式（預設："text"）
}
```

## 輸出格式

### JSON 格式（slog 標準）
當 `Type: "json"` 時，日誌以 slog 相容的結構化格式輸出：

```json
{"timestamp":"2024/01/15 14:30:25.123456","level":"INFO","message":"應用程式已啟動","data":null}
{"timestamp":"2024/01/15 14:30:25.123457","level":"ERROR","message":"資料庫連線失敗","data":["連線逾時","5 秒後重試"]}
```

優點：
- 與 Go 標準 `log/slog` 套件相容
- 機器可讀的結構化日誌記錄
- 易於與日誌聚合工具整合
- 所有日誌層級保持一致的 JSON 架構

### 文字格式（樹狀結構）
當 `Type: "text"` 時，日誌以人類可讀的樹狀格式顯示：

```
2024/01/15 14:30:25.123456 應用程式已啟動
2024/01/15 14:30:25.123457 [ERROR] 資料庫連線失敗
2024/01/15 14:30:25.123457 ├── 連線逾時
2024/01/15 14:30:25.123457 └── 5 秒後重試
```

優點：
- 人類友善的視覺化呈現
- 清晰的階層訊息結構
- 提升除錯時的可讀性

## 日誌層級說明

### Debug 和 Trace
記錄到 `debug.log`
```go
logger.Debug("變數值", "x = 10", "y = 20")
logger.Trace("函式呼叫", "開始處理使用者請求")
```

### Info、Notice、Warning
記錄到 `output.log`
```go
logger.Info("應用程式已啟動")           // 無前綴
logger.Notice("設定檔已重新載入") // [NOTICE] 前綴
logger.Warn("記憶體使用量過高")         // [WARNING] 前綴
```

### Error、Fatal、Critical
記錄到 `error.log`
```go
logger.Error(err, "重試第 3 次")        // [ERROR] 前綴
logger.Fatal(err, "無法啟動服務")   // [FATAL] 前綴
logger.Critical(err, "系統當機")        // [CRITICAL] 前綴
```

## 核心功能

### 日誌記錄器管理

- **New** - 建立新的日誌記錄器實例
  ```go
  logger, err := goLogger.New(config)
  ```
  - 初始化日誌目錄，確保路徑存在
  - 開啟三個日誌檔案：debug.log、output.log、error.log
  - 為每個層級設定日誌處理器
  - 檢查現有檔案大小，必要時執行輪替

- **Close** - 安全關閉日誌記錄器
  ```go
  err := logger.Close()
  ```
  - 關閉所有開啟的檔案控制代碼
  - 標記日誌記錄器為已關閉
  - 確保無資源洩漏

- **Flush** - 強制寫入快取
  ```go
  err := logger.Flush()
  ```
  - 將所有快取的日誌內容寫入磁碟
  - 確保重要日誌不會遺失

### 檔案輪替機制

#### 自動輪替
- 每次日誌寫入前檢查檔案大小
- 超過 `MaxSize` 限制時自動輪替
- 備份檔案命名格式：`filename.YYYYMMDD_HHMMSS`

#### 備份管理
- 保留最新的 `MaxBackup` 個備份檔案
- 自動刪除過期的舊備份
- 按修改時間排序，保留最新的檔案

### 並行安全機制

#### 讀寫鎖保護
- 使用 `sync.RWMutex` 保護關鍵區段
- 寫入操作取得寫入鎖，確保原子性
- 讀取操作使用讀取鎖，提升並行效能

## 使用範例

### JSON 格式的基本日誌記錄
```go
logger, _ := goLogger.New(&goLogger.Log{
  Path:    "./logs",
  Stdout:  true,
  MaxSize: 1024 * 1024, // 1MB
  Type:    "json",      // slog 標準格式
})
defer logger.Close()

logger.Info("應用程式已啟動")
logger.Debug("載入設定檔", "config.json")
logger.Warn("記憶體使用量", "85%")
```

### 文字格式的基本日誌記錄
```go
logger, _ := goLogger.New(&goLogger.Log{
  Path:    "./logs",
  Stdout:  true,
  MaxSize: 1024 * 1024, // 1MB
  Type:    "text",      // 樹狀結構格式
})
defer logger.Close()

logger.Info("應用程式已啟動")
logger.Debug("載入設定檔", "config.json")
logger.Warn("記憶體使用量", "85%")
```

### 錯誤處理
```go
if err := connectDatabase(); err != nil {
  logger.Error(err, "資料庫連線失敗", "重試中...")
  return logger.Fatal(err, "無法建立資料庫連線")
}
```

### 並行環境
```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
    defer wg.Done()
    logger.Info(fmt.Sprintf("Goroutine %d 執行中", id))
  }(i)
}

wg.Wait()
logger.Flush() // 確保所有日誌都已寫入
```

## 授權條款

此原始碼專案採用 [MIT](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) 授權條款。

## 作者

<img src="https://avatars.githubusercontent.com/u/25631760" align="left" width="96" height="96" style="margin-right: 0.5rem;">

<h4 style="padding-top: 0">邱敬幃 Pardn Chiu</h4>

<a href="mailto:dev@pardn.io" target="_blank">
    <img src="https://pardn.io/image/email.svg" width="48" height="48">
</a> <a href="https://linkedin.com/in/pardnchiu" target="_blank">
    <img src="https://pardn.io/image/linkedin.svg" width="48" height="48">
</a>

***

©️ 2025 [邱敬幃 Pardn Chiu](https://pardn.io)