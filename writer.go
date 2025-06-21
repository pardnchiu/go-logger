package logger

import (
	"fmt"
	"log"
	"strings"
)

func (l *Logger) writeToLog(target *log.Logger, level string, filename string, messages ...string) {
	level = strings.ToUpper(level)
	isValid := map[string]bool{
		logDebug:    true,
		logTrace:    true,
		logInfo:     true,
		logNotice:   true,
		logWarning:  true,
		logError:    true,
		logFatal:    true,
		logCritical: true,
	}[level]

	if !isValid {
		return
	}

	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.IsClose || len(messages) == 0 {
		return
	}

	if err := l.checkAndRotate(filename); err != nil {
		fmt.Printf("Warning: log rotation failed: %v\n", err)
	}

	prefix := ""
	if level != logInfo {
		prefix = fmt.Sprintf("[%s] ", level)
	}

	for i, msg := range messages {
		switch {
		case i == 0:
			target.Printf("%s%s", prefix, msg)
		case i == len(messages)-1:
			target.Printf("└── %s", msg)
		default:
			target.Printf("├── %s", msg)
		}
	}
}

func (l *Logger) Debug(messages ...string) {
	l.writeToLog(l.DebugHandler, logDebug, defaultDebugName, messages...)
}

func (l *Logger) Trace(messages ...string) {
	l.writeToLog(l.DebugHandler, logTrace, defaultDebugName, messages...)
}

func (l *Logger) Info(messages ...string) {
	l.writeToLog(l.OutputHandler, logInfo, defaultOutputName, messages...)
}

func (l *Logger) Notice(messages ...string) {
	l.writeToLog(l.OutputHandler, logNotice, defaultOutputName, messages...)
}

func (l *Logger) Warn(messages ...string) {
	l.writeToLog(l.OutputHandler, logWarning, defaultOutputName, messages...)
}

func (l *Logger) Error(err error, messages ...string) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logError, defaultErrorName, messages...)
	return fmt.Errorf("%s", strings.Join(messages, " "))
}

func (l *Logger) Fatal(err error, messages ...string) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logFatal, defaultErrorName, messages...)
	return fmt.Errorf("%s", strings.Join(messages, " "))
}

func (l *Logger) Critical(err error, messages ...string) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logCritical, defaultErrorName, messages...)
	return fmt.Errorf("%s", strings.Join(messages, " "))
}
