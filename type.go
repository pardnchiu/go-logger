package logger

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultDebugName  = "debug.log"
	defaultOutputName = "output.log"
	defaultErrorName  = "error.log"
	logDebug          = "DEBUG"
	logTrace          = "TRACE"
	logInfo           = "INFO"
	logNotice         = "NOTICE"
	logWarning        = "WARNING"
	logError          = "ERROR"
	logFatal          = "FATAL"
	logCritical       = "CRITICAL"
)

type Log struct {
	Path      string `json:"path,omitempty"`        // 日誌檔案路徑，預設 `./logs`
	Stdout    bool   `json:"stdout,omitempty"`      // 是否輸出到標準輸出，預設 false
	MaxSize   int64  `json:"max_size,omitempty"`    // 日誌檔案最大大小（位元組），預設 16 * 1024 * 1024
	MaxBackup int    `json:"max_backups,omitempty"` // 新增：最大備份檔案數量，預設 5
}

type Logger struct {
	Config        *Log
	DebugHandler  *log.Logger
	OutputHandler *log.Logger
	ErrorHandler  *log.Logger
	File          map[string]*os.File
	Mutex         sync.RWMutex
	IsClose       bool
}

type backupFile struct {
	path    string
	modTime time.Time
}
