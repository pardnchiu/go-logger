# Go Logger (Golang)

> 一個專為 Golang 設計的日誌套件，具備自動輪替、多層級日誌分類和檔案管理功能，並提供完整的錯誤處理機制。<br>
> 主要用於 `pardnchiu/go-*` 套件系列

[![license](https://img.shields.io/github/license/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) 
[![version](https://img.shields.io/github/v/tag/pardnchiu/go-logger)](https://github.com/pardnchiu/go-logger/releases) 
[![readme](https://img.shields.io/badge/readme-English-blue)](https://github.com/pardnchiu/go-logger/blob/main/README.md) 

## 功能特色

- **多層級日誌分類**：支援 DEBUG、TRACE、INFO、NOTICE、WARNING、ERROR、FATAL、CRITICAL 層級
- **自動檔案輪替**：當檔案大小超過限制時自動建立備份並開始新檔案
- **備份檔案管理**：自動清理過期備份，維護可設定的備份數量
- **併發安全**：執行緒安全的日誌寫入，支援高併發環境
- **多重輸出目標**：同時輸出至檔案和標準輸出
- **樹狀結構訊息**：多行訊息以樹狀結構顯示，提升可讀性
- **記憶體高效**：基於互斥鎖的安全寫入，防止資料競爭

## 使用方式

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
    Stdout:    true,                  // 同時輸出至終端
    MaxSize:   16 * 1024 * 1024,      // 16MB 檔案大小限制
    MaxBackup: 5,                     // 保留 5 個備份檔案
  }
  
  // 初始化記錄器
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
  logger.Error(err, "處理錯誤時的額外訊息")
  logger.Fatal(err, "嚴重錯誤")
  logger.Critical(err, "系統關鍵錯誤")
  
  // 手動清空快取
  logger.Flush()
}
```

### 設定詳情

```go
type Log struct {
  Path      string // 日誌檔案目錄路徑（預設：./logs）
  Stdout    bool   // 是否輸出至標準輸出（預設：false）
  MaxSize   int64  // 最大日誌檔案大小（位元組）（預設：16MB）
  MaxBackup int    // 最大備份檔案數量（預設：5）
}
```

## 日誌層級說明

### Debug 和 Trace
記錄至 `debug.log`
```go
logger.Debug("變數值", "x = 10", "y = 20")
logger.Trace("函式呼叫", "開始處理使用者請求")
```

### Info、Notice、Warning
記錄至 `output.log`
```go
logger.Info("應用程式啟動")           // 無前綴
logger.Notice("設定檔重新載入")        // [NOTICE] 前綴
logger.Warn("記憶體使用量過高")        // [WARNING] 前綴
```

### Error、Fatal、Critical
記錄至 `error.log`
```go
logger.Error(err, "重試第 3 次")       // [ERROR] 前綴
logger.Fatal(err, "無法啟動服務")      // [FATAL] 前綴
logger.Critical(err, "系統當機")       // [CRITICAL] 前綴
```

## 核心功能

### 記錄器管理

- **New** - 建立新記錄器實例
  ```go
  logger, err := goLogger.New(config)
  ```
  - 初始化日誌目錄，確保路徑存在
  - 開啟三個日誌檔案：debug.log、output.log、error.log
  - 為每個層級設定日誌處理器
  - 檢查現有檔案大小，必要時執行輪替

- **Close** - 安全關閉記錄器
  ```go
  err := logger.Close()
  ```
  - 關閉所有開啟的檔案控制代碼
  - 標記記錄器為已關閉
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
- 保留最新的 `MaxBackup` 備份檔案
- 自動刪除過期的舊備份
- 依修改時間排序，保留最新檔案

### 併發安全機制

#### 讀寫鎖保護
- 使用 `sync.RWMutex` 保護關鍵區段
- 寫入操作取得寫入鎖，確保原子性
- 讀取操作使用讀取鎖，提升併發效能

## 訊息格式化

### 單行訊息
```go
logger.Info("單一訊息")
```
輸出：
```
2024/01/15 14:30:25.123456 單一訊息
```

### 多行樹狀結構
```go
logger.Error(err, "主要錯誤", "詳細資訊", "額外備註")
```
輸出：
```
2024/01/15 14:30:25.123456 [ERROR] 主要錯誤
2024/01/15 14:30:25.123456 ├── 詳細資訊
2024/01/15 14:30:25.123456 └── 額外備註
```

## 使用範例

### 基本日誌記錄
```go
logger, _ := goLogger.New(&goLogger.Log{
  Path:    "./logs",
  Stdout:  true,
  MaxSize: 1024 * 1024, // 1MB
})
defer logger.Close()

logger.Info("應用程式啟動")
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

### 併發環境
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

本原始碼專案採用 [MIT](https://github.com/pardnchiu/go-logger/blob/main/LICENSE) 授權條款。

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