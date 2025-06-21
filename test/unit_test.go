package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	goLogger "github.com/pardnchiu/go-logger"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		config   *goLogger.Log
		expected *goLogger.Log
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
			expected: &goLogger.Log{
				Path:      "./logs",
				Stdout:    false,
				MaxSize:   16 * 1024 * 1024,
				MaxBackup: 5,
			},
		},
		{
			name: "empty path uses default",
			config: &goLogger.Log{
				Path:      "",
				Stdout:    true,
				MaxSize:   1024,
				MaxBackup: 3,
			},
			expected: &goLogger.Log{
				Path:      "./logs",
				Stdout:    true,
				MaxSize:   1024,
				MaxBackup: 3,
			},
		},
		{
			name: "zero values use defaults",
			config: &goLogger.Log{
				Path:      "./test_logs",
				Stdout:    false,
				MaxSize:   0,
				MaxBackup: 0,
			},
			expected: &goLogger.Log{
				Path:      "./test_logs",
				Stdout:    false,
				MaxSize:   16 * 1024 * 1024,
				MaxBackup: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			if tt.expected.Path != "" {
				os.RemoveAll(tt.expected.Path)
			}

			logger, err := goLogger.New(tt.config)
			if err != nil {
				t.Fatalf("newLogger() error = %v", err)
			}
			defer logger.Close()

			if logger.Config.Path != tt.expected.Path {
				t.Errorf("Path = %v, want %v", logger.Config.Path, tt.expected.Path)
			}
			if logger.Config.Stdout != tt.expected.Stdout {
				t.Errorf("Stdout = %v, want %v", logger.Config.Stdout, tt.expected.Stdout)
			}
			if logger.Config.MaxSize != tt.expected.MaxSize {
				t.Errorf("MaxSize = %v, want %v", logger.Config.MaxSize, tt.expected.MaxSize)
			}
			if logger.Config.MaxBackup != tt.expected.MaxBackup {
				t.Errorf("MaxBackup = %v, want %v", logger.Config.MaxBackup, tt.expected.MaxBackup)
			}

			// Clean up after test
			os.RemoveAll(tt.expected.Path)
		})
	}
}

func TestLoggerInit(t *testing.T) {
	testDir := "./test_logs_init"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	// Check if all log files are created
	expectedFiles := []string{"debug.log", "output.log", "error.log"}
	for _, filename := range expectedFiles {
		fullPath := filepath.Join(testDir, filename)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", fullPath)
		}
	}

	// Check if handlers are initialized
	if logger.DebugHandler == nil {
		t.Error("debugHandler was not initialized")
	}
	if logger.OutputHandler == nil {
		t.Error("outputHandler was not initialized")
	}
	if logger.ErrorHandler == nil {
		t.Error("errorHandler was not initialized")
	}
}

func TestLoggerLoggingMethods(t *testing.T) {
	testDir := "./test_logs_methods"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024 * 1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	// Test debug and trace methods
	logger.Debug("Debug message 1", "Debug message 2")
	logger.Trace("Trace message")

	// Test info, notice, warning methods
	logger.Info("Info message")
	logger.Notice("Notice message")
	logger.Warn("Warning message")

	// Test error, fatal, critical methods
	testErr := fmt.Errorf("test error")
	logger.Error(testErr, "Error message")
	logger.Fatal(testErr, "Fatal message")
	logger.Critical(testErr, "Critical message")

	// Flush to ensure all logs are written
	logger.Flush()

	// Verify debug log content
	debugContent := readLogFile(t, filepath.Join(testDir, "debug.log"))
	if !strings.Contains(debugContent, "[DEBUG]") {
		t.Error("Debug log should contain [DEBUG] prefix")
	}
	if !strings.Contains(debugContent, "[TRACE]") {
		t.Error("Debug log should contain [TRACE] prefix")
	}

	// Verify output log content
	outputContent := readLogFile(t, filepath.Join(testDir, "output.log"))
	if !strings.Contains(outputContent, "Info message") {
		t.Error("Output log should contain info message without prefix")
	}
	if !strings.Contains(outputContent, "[NOTICE]") {
		t.Error("Output log should contain [NOTICE] prefix")
	}
	if !strings.Contains(outputContent, "[WARNING]") {
		t.Error("Output log should contain [WARNING] prefix")
	}

	// Verify error log content
	errorContent := readLogFile(t, filepath.Join(testDir, "error.log"))
	if !strings.Contains(errorContent, "[ERROR]") {
		t.Error("Error log should contain [ERROR] prefix")
	}
	if !strings.Contains(errorContent, "[FATAL]") {
		t.Error("Error log should contain [FATAL] prefix")
	}
	if !strings.Contains(errorContent, "[CRITICAL]") {
		t.Error("Error log should contain [CRITICAL] prefix")
	}
}

func TestLogRotation(t *testing.T) {
	testDir := "./test_logs_rotation"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   100, // Very small size to trigger rotation
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	// Write enough data to trigger rotation
	for i := 0; i < 50; i++ {
		logger.Info(fmt.Sprintf("This is a long message to trigger rotation - message %d", i))
	}

	// Check if backup files are created
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	backupPattern := regexp.MustCompile(`output\.log\.\d{8}_\d{6}`)
	backupCount := 0
	for _, file := range files {
		if backupPattern.MatchString(file.Name()) {
			backupCount++
		}
	}

	if backupCount == 0 {
		t.Error("Expected at least one backup file to be created")
	}
}

func TestLogCleanup(t *testing.T) {
	testDir := "./test_logs_cleanup"
	defer os.RemoveAll(testDir)

	// Create test directory
	os.MkdirAll(testDir, 0755)

	// Create old backup files
	baseFile := "output.log"
	testFiles := []string{
		"output.log.20230101_120000",
		"output.log.20230102_120000",
		"output.log.20230103_120000",
		"output.log.20230104_120000",
		"output.log.20230105_120000",
		"output.log.20230106_120000", // This should be removed
	}

	for i, filename := range testFiles {
		fullPath := filepath.Join(testDir, filename)
		file, _ := os.Create(fullPath)
		file.Close()

		// Set different modification times
		modTime := time.Now().AddDate(0, 0, -len(testFiles)+i)
		os.Chtimes(fullPath, modTime, modTime)
	}

	config := &goLogger.Log{
		Path:      testDir,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}

	// Test cleanup
	err = logger.Cleanup(filepath.Join(testDir, baseFile))
	if err != nil {
		t.Fatalf("cleanup() error = %v", err)
	}

	// Count remaining backup files
	files, _ := os.ReadDir(testDir)
	backupPattern := regexp.MustCompile(`^output\.log\.\d{8}_\d{6}$`)
	backupCount := 0
	for _, file := range files {
		if backupPattern.MatchString(file.Name()) {
			backupCount++
		}
	}

	if backupCount != config.MaxBackup {
		t.Errorf("Expected %d backup files, got %d", config.MaxBackup, backupCount)
	}
}

func TestConcurrentLogging(t *testing.T) {
	testDir := "./test_logs_concurrent"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024 * 1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	var wg sync.WaitGroup
	numRoutines := 10
	messagesPerRoutine := 100

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < messagesPerRoutine; j++ {
				logger.Info(fmt.Sprintf("Routine %d - Message %d", routineID, j))
				logger.Debug(fmt.Sprintf("Debug from routine %d - Message %d", routineID, j))
				logger.Error(nil, fmt.Sprintf("Error from routine %d - Message %d", routineID, j))
			}
		}(i)
	}

	wg.Wait()

	// Verify no race conditions occurred by checking if files can be read
	logger.Flush()
	outputContent := readLogFile(t, filepath.Join(testDir, "output.log"))
	debugContent := readLogFile(t, filepath.Join(testDir, "debug.log"))
	errorContent := readLogFile(t, filepath.Join(testDir, "error.log"))

	if len(outputContent) == 0 {
		t.Error("Output log should contain concurrent messages")
	}
	if len(debugContent) == 0 {
		t.Error("Debug log should contain concurrent messages")
	}
	if len(errorContent) == 0 {
		t.Error("Error log should contain concurrent messages")
	}
}

func TestLoggerClose(t *testing.T) {
	testDir := "./test_logs_close"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}

	// Write some data
	logger.Info("Test message before close")

	// Close the logger
	err = logger.Close()
	if err != nil {
		t.Errorf("close() error = %v", err)
	}

	// Verify logger is closed
	if !logger.IsClose {
		t.Error("Logger should be marked as closed")
	}

	// Verify subsequent operations don't crash
	logger.Info("This should not cause panic")

	// Verify closing again doesn't error
	err = logger.Close()
	if err != nil {
		t.Errorf("Second close() should not error, got %v", err)
	}
}

func TestLoggerFlush(t *testing.T) {
	testDir := "./test_logs_flush"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	// Write some data
	logger.Info("Test message")

	// Test flush
	err = logger.Flush()
	if err != nil {
		t.Errorf("flush() error = %v", err)
	}

	// Close logger and test flush on closed logger
	logger.Close()
	err = logger.Flush()
	if err == nil {
		t.Error("flush() on closed logger should return error")
	}
}

func TestEmptyMessages(t *testing.T) {
	testDir := "./test_logs_empty"
	defer os.RemoveAll(testDir)

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024,
		MaxBackup: 3,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("newLogger() error = %v", err)
	}
	defer logger.Close()

	// Test with no messages (should not crash)
	logger.Info()
	logger.Debug()
	logger.Error(nil)

	// Verify this doesn't cause issues
	logger.Flush()
}

// Helper function to read log file content
func readLogFile(t *testing.T, filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read log file %s: %v", filePath, err)
	}
	return string(content)
}
