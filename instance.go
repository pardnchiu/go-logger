package goLogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

func New(config *Log) (*Logger, error) {
	if config == nil {
		config = &Log{
			Path:      "./logs",
			Stdout:    false,
			MaxSize:   16 * 1024 * 1024,
			MaxBackup: 5,
		}
	}
	if config.Path == "" {
		config.Path = "./logs"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 16 * 1024 * 1024
	}
	if config.MaxBackup == 0 {
		config.MaxBackup = 5
	}
	if config.Type == "" {
		config.Type = "text"
	}

	if err := os.MkdirAll(config.Path, 0755); err != nil {
		return nil, fmt.Errorf("Failed to create: %w", err)
	}

	logger := &Logger{
		Config: config,
		File:   make(map[string]*os.File),
	}

	if err := logger.init(0644); err != nil {
		logger.Close()
		return nil, err
	}

	logger.startRotateTimer()

	return logger, nil
}

func (l *Logger) init(mode os.FileMode) error {
	files := []string{defaultDebugName, defaultOutputName, defaultErrorName}

	for _, filename := range files {
		file, err := l.open(filename, mode)
		if err != nil {
			return err
		}
		l.File[filename] = file
	}

	return l.initHandler()
}

func (l *Logger) initHandler() error {
	flags := log.LstdFlags | log.Lmicroseconds

	var debugWriters []io.Writer = []io.Writer{l.File[defaultDebugName]}
	var outputWriters []io.Writer = []io.Writer{l.File[defaultOutputName]}
	var errorWriters []io.Writer = []io.Writer{l.File[defaultErrorName]}

	if l.Config.Stdout {
		debugWriters = append(debugWriters, os.Stdout)
		outputWriters = append(outputWriters, os.Stdout)
		errorWriters = append(errorWriters, os.Stderr)
	}

	l.DebugHandler = log.New(io.MultiWriter(debugWriters...), "", flags)
	l.OutputHandler = log.New(io.MultiWriter(outputWriters...), "", flags)
	l.ErrorHandler = log.New(io.MultiWriter(errorWriters...), "", flags)

	return nil
}

func (l *Logger) open(filename string, mode os.FileMode) (*os.File, error) {
	fullPath := filepath.Join(l.Config.Path, filename)

	if info, err := os.Stat(fullPath); err == nil {
		// * file exists
		if info.Size() > l.Config.MaxSize {
			// * size exceeds max size
			if err := l.rotate(fullPath); err != nil {
				// * failed to rotate
				return nil, fmt.Errorf("Failed to rotate %s: %w", filename, err)
			}
		}
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %w", filename, err)
	}
	return file, nil
}

func (l *Logger) rotate(path string) error {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", path, timestamp)

	if err := os.Rename(path, backupPath); err != nil {
		// * failed to rename old log
		return fmt.Errorf("Failed to rotate: %w", err)
	}

	if err := l.Cleanup(path); err != nil {
		fmt.Printf("Failed to clean: %v", err)
	}

	return nil
}

func (l *Logger) Cleanup(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Failed to read: %w", err)
	}

	backupPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(base) + `\.\d{8}_\d{6}$`)

	var backupFiles []backupFile
	for _, file := range files {
		name := file.Name()
		// * filename.YYYYMMDD_HHMMSS
		if backupPattern.MatchString(name) {
			info, err := file.Info()
			if err != nil {
				continue
			}

			backupFiles = append(backupFiles, backupFile{
				path:    filepath.Join(dir, name),
				modTime: info.ModTime(),
			})
		}
	}

	if len(backupFiles) > l.Config.MaxBackup {
		sort.Slice(backupFiles, func(i, j int) bool {
			return backupFiles[i].modTime.After(backupFiles[j].modTime)
		})

		for i := l.Config.MaxBackup; i < len(backupFiles); i++ {
			if err := os.Remove(backupFiles[i].path); err != nil {
				return fmt.Errorf("Failed to remove %s: %w", backupFiles[i].path, err)
			}
		}
	}

	return nil
}

func (l *Logger) startRotateTimer() {
	l.stopTimer = make(chan struct{})
	l.timer = time.NewTimer(1 * time.Hour)

	go func() {
		for {
			select {
			case <-l.timer.C:
				l.checkAndRotate(defaultDebugName)
				l.checkAndRotate(defaultOutputName)
				l.checkAndRotate(defaultErrorName)
				l.timer.Reset(1 * time.Hour)
			case <-l.stopTimer:
				if l.timer != nil {
					l.timer.Stop()
				}
				return
			}
		}
	}()
}

func (l *Logger) checkAndRotate(filename string) error {
	oldFile, isExist := l.File[filename]
	if !isExist {
		return fmt.Errorf("Failed to read: %s", filename)
	}

	stat, err := oldFile.Stat()
	if err != nil {
		return fmt.Errorf("Failed to get stats: %w", err)
	}

	if stat.Size() > l.Config.MaxSize {
		oldFile.Close()

		path := filepath.Join(l.Config.Path, filename)
		if err := l.rotate(path); err != nil {
			return fmt.Errorf("Failed to rotate %s: %w", filename, err)
		}

		newFile, err := l.open(filename, 0644)
		if err != nil {
			return fmt.Errorf("Failed to reopen %s: %w", filename, err)
		}

		l.File[filename] = newFile

		if err := l.initHandler(); err != nil {
			return fmt.Errorf("Failed to re-init: %w", err)
		}
	}

	return nil
}

func (l *Logger) Close() error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.IsClose {
		return nil
	}

	l.IsClose = true

	if l.stopTimer != nil {
		close(l.stopTimer)
	}

	var errs []error

	for filename, file := range l.File {
		if err := file.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing %s: %w", filename, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing log files: %v", errs)
	}

	return nil
}

func (l *Logger) Flush() error {
	l.Mutex.RLock()
	defer l.Mutex.RUnlock()

	if l.IsClose {
		return fmt.Errorf("logger is closed")
	}

	var errs []error
	for filename, file := range l.File {
		if err := file.Sync(); err != nil {
			errs = append(errs, fmt.Errorf("flushing %s: %w", filename, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors flushing log files: %v", errs)
	}

	return nil
}
