package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	goLogger "github.com/pardnchiu/go-logger" // Adjust the import path to your logger package
)

func createTestLogger(t *testing.T, logType string) (*goLogger.Logger, string) {
	testDir := fmt.Sprintf("./test_writer_%s_%d", logType, time.Now().UnixNano())

	config := &goLogger.Log{
		Path:      testDir,
		Stdout:    false,
		MaxSize:   1024,
		MaxBackup: 3,
		Type:      logType,
	}

	logger, err := goLogger.New(config)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	return logger, testDir
}

func readLogContent(t *testing.T, filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read log file %s: %v", filepath, err)
	}
	return string(content)
}

func TestDebugLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Debug("Debug message")
	logger.Debug("Debug with", "multiple", "arguments")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "debug.log"))

	if !strings.Contains(content, "Debug message") {
		t.Error("Debug log should contain debug message")
	}
	if !strings.Contains(content, `"level":"DEBUG"`) {
		t.Error("JSON debug log should contain DEBUG level")
	}
}

func TestTraceLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Trace("Trace message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "debug.log"))

	if !strings.Contains(content, "Trace message") {
		t.Error("Trace log should contain trace message")
	}
	if !strings.Contains(content, `"level":"TRACE"`) {
		t.Error("JSON trace log should contain TRACE level")
	}
}

func TestInfoLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Info("Info message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))

	if !strings.Contains(content, "Info message") {
		t.Error("Info log should contain info message")
	}
	if !strings.Contains(content, `"level":"INFO"`) {
		t.Error("JSON info log should contain INFO level")
	}
}

func TestNoticeLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Notice("Notice message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))

	if !strings.Contains(content, "Notice message") {
		t.Error("Notice log should contain notice message")
	}
	if !strings.Contains(content, `"level":"NOTICE"`) {
		t.Error("JSON notice log should contain NOTICE level")
	}
}

func TestWarnLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Warn("Warning message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))

	if !strings.Contains(content, "Warning message") {
		t.Error("Warn log should contain warning message")
	}
	if !strings.Contains(content, `"level":"WARN"`) {
		t.Error("JSON warn log should contain WARN level")
	}
}

func TestErrorLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	testError := fmt.Errorf("test error")
	returnedError := logger.Error(testError, "Error message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "error.log"))

	if !strings.Contains(content, "Error message") {
		t.Error("Error log should contain error message")
	}
	if !strings.Contains(content, "test error") {
		t.Error("Error log should contain error details")
	}
	if !strings.Contains(content, `"level":"ERROR"`) {
		t.Error("JSON error log should contain ERROR level")
	}
	if returnedError == nil {
		t.Error("Error method should return an error")
	}
}

func TestErrorLoggingWithNilError(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	returnedError := logger.Error(nil, "Error message without error object")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "error.log"))

	if !strings.Contains(content, "Error message without error object") {
		t.Error("Error log should contain error message even with nil error")
	}
	if returnedError == nil {
		t.Error("Error method should return an error even with nil input error")
	}
}

func TestFatalLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	testError := fmt.Errorf("fatal error")
	returnedError := logger.Fatal(testError, "Fatal message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "error.log"))

	if !strings.Contains(content, "Fatal message") {
		t.Error("Fatal log should contain fatal message")
	}
	if !strings.Contains(content, "fatal error") {
		t.Error("Fatal log should contain error details")
	}
	if !strings.Contains(content, `"level":"FATAL"`) {
		t.Error("JSON fatal log should contain FATAL level")
	}
	if returnedError == nil {
		t.Error("Fatal method should return an error")
	}
}

func TestCriticalLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	testError := fmt.Errorf("critical error")
	returnedError := logger.Critical(testError, "Critical message")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "error.log"))

	if !strings.Contains(content, "Critical message") {
		t.Error("Critical log should contain critical message")
	}
	if !strings.Contains(content, "critical error") {
		t.Error("Critical log should contain error details")
	}
	if !strings.Contains(content, `"level":"CRITICAL"`) {
		t.Error("JSON critical log should contain CRITICAL level")
	}
	if returnedError == nil {
		t.Error("Critical method should return an error")
	}
}

func TestTextFormatLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "text")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Debug("Debug text message")
	logger.Info("Info text message")
	logger.Error(nil, "Error text message")
	logger.Flush()

	debugContent := readLogContent(t, filepath.Join(testDir, "debug.log"))
	outputContent := readLogContent(t, filepath.Join(testDir, "output.log"))
	errorContent := readLogContent(t, filepath.Join(testDir, "error.log"))

	// Check text format (should not contain JSON)
	if strings.Contains(debugContent, `"level"`) {
		t.Error("Text format should not contain JSON level field")
	}
	if strings.Contains(outputContent, `"level"`) {
		t.Error("Text format should not contain JSON level field")
	}
	if strings.Contains(errorContent, `"level"`) {
		t.Error("Text format should not contain JSON level field")
	}

	// Check text format prefixes
	if !strings.Contains(debugContent, "[DEBUG]") {
		t.Error("Text debug log should contain [DEBUG] prefix")
	}
	if strings.Contains(outputContent, "[INFO]") {
		t.Error("Text info log should not contain [INFO] prefix")
	}
	if !strings.Contains(errorContent, "[ERROR]") {
		t.Error("Text error log should contain [ERROR] prefix")
	}
}

func TestMultipleArgumentsTextFormat(t *testing.T) {
	logger, testDir := createTestLogger(t, "text")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Info("Main message", "Second argument", "Third argument")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))

	if !strings.Contains(content, "Main message") {
		t.Error("Text log should contain main message")
	}
	if !strings.Contains(content, "├── Second argument") {
		t.Error("Text log should contain tree structure for middle arguments")
	}
	if !strings.Contains(content, "└── Third argument") {
		t.Error("Text log should contain tree structure for last argument")
	}
}

func TestMultipleArgumentsJSONFormat(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	logger.Info("Main message", "Second argument", "Third argument")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))

	var logEntry map[string]interface{}
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON log: %v", err)
		}

		if logEntry["msg"] != "Main message" {
			t.Error("JSON log should contain main message in msg field")
		}
		if logEntry["msg1"] != "Second argument" {
			t.Error("JSON log should contain second argument in msg1 field")
		}
		if logEntry["msg2"] != "Third argument" {
			t.Error("JSON log should contain third argument in msg2 field")
		}
	}
}

func TestEmptyMessages(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	// Should not log anything with empty messages
	logger.Info()
	logger.Debug()
	logger.Error(nil)
	logger.Flush()

	// Check that no content was written
	outputContent := readLogContent(t, filepath.Join(testDir, "output.log"))
	debugContent := readLogContent(t, filepath.Join(testDir, "debug.log"))
	errorContent := readLogContent(t, filepath.Join(testDir, "error.log"))

	if strings.TrimSpace(outputContent) != "" {
		t.Error("Empty message should not write to output log")
	}
	if strings.TrimSpace(debugContent) != "" {
		t.Error("Empty message should not write to debug log")
	}
	if strings.TrimSpace(errorContent) != "" {
		t.Error("Empty message should not write to error log")
	}
}

func TestClosedLogger(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)

	// Close the logger
	logger.Close()

	// Try to log after closing
	logger.Info("This should not be logged")
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))
	if strings.Contains(content, "This should not be logged") {
		t.Error("Closed logger should not log messages")
	}
}

func TestConcurrentLogging(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	var wg sync.WaitGroup
	numGoroutines := 10
	messagesPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info(fmt.Sprintf("Goroutine %d message %d", id, j))
			}
		}(i)
	}

	wg.Wait()
	logger.Flush()

	content := readLogContent(t, filepath.Join(testDir, "output.log"))
	lines := strings.Split(strings.TrimSpace(content), "\n")

	// Should have all messages logged
	expectedMessages := numGoroutines * messagesPerGoroutine
	if len(lines) != expectedMessages {
		t.Errorf("Expected %d log lines, got %d", expectedMessages, len(lines))
	}
}

func TestLogRotationTrigger(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	// Set very small max size to trigger rotation
	logger.Config.MaxSize = 10

	// Log enough data to trigger rotation
	for i := 0; i < 100; i++ {
		logger.Info(fmt.Sprintf("This is a long message to trigger log rotation %d", i))
	}
	logger.Flush()

	// Check that rotation was attempted (files should exist)
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	if len(files) < 3 { // Should have at least debug.log, output.log, error.log
		t.Error("Log rotation should maintain log files")
	}
}

func TestNilErrorInAllErrorMethods(t *testing.T) {
	logger, testDir := createTestLogger(t, "json")
	defer os.RemoveAll(testDir)
	defer logger.Close()

	// Test all error methods with nil error
	errorResult := logger.Error(nil, "Error with nil")
	fatalResult := logger.Fatal(nil, "Fatal with nil")
	criticalResult := logger.Critical(nil, "Critical with nil")
	logger.Flush()

	// All should return non-nil errors
	if errorResult == nil {
		t.Error("Error method should return error even with nil input")
	}
	if fatalResult == nil {
		t.Error("Fatal method should return error even with nil input")
	}
	if criticalResult == nil {
		t.Error("Critical method should return error even with nil input")
	}

	content := readLogContent(t, filepath.Join(testDir, "error.log"))
	if !strings.Contains(content, "Error with nil") {
		t.Error("Error log should contain error message")
	}
	if !strings.Contains(content, "Fatal with nil") {
		t.Error("Error log should contain fatal message")
	}
	if !strings.Contains(content, "Critical with nil") {
		t.Error("Error log should contain critical message")
	}
}
