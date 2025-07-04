package goLogger

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
)

func (l *Logger) writeToLog(target *log.Logger, level string, filename string, messages ...any) {
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

	if l.Config.Type == "json" {
		jsonLogger := slog.New(slog.NewJSONHandler(target.Writer(), &slog.HandlerOptions{
			Level: slog.LevelDebug, // 確保 DEBUG 層級會被輸出
		}))

		msg := fmt.Sprintf("%v", messages[0])
		remaining := messages[1:]
		attrs := make([]any, len(remaining))
		for i, m := range remaining {
			attrs[i] = slog.String(fmt.Sprintf("msg%d", i+1), fmt.Sprintf("%v", m))
		}

		switch level {
		case logDebug:
			jsonLogger.Debug(msg, attrs...)
		case logTrace:
			jsonLogger.Info(msg, append(attrs, slog.String("level", "TRACE"))...)
		case logInfo:
			jsonLogger.Info(msg, attrs...)
		case logNotice:
			jsonLogger.Info(msg, append(attrs, slog.String("level", "NOTICE"))...)
		case logWarning:
			jsonLogger.Warn(msg, attrs...)
		case logError:
			jsonLogger.Error(msg, attrs...)
		case logFatal:
			jsonLogger.Error(msg, append(attrs, slog.String("level", "FATAL"))...)
		case logCritical:
			jsonLogger.Error(msg, append(attrs, slog.String("level", "CRITICAL"))...)
		}
		return
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

func (l *Logger) Debug(messages ...any) {
	l.writeToLog(l.DebugHandler, logDebug, defaultDebugName, messages...)
}

func (l *Logger) Trace(messages ...any) {
	l.writeToLog(l.DebugHandler, logTrace, defaultDebugName, messages...)
}

func (l *Logger) Info(messages ...any) {
	l.writeToLog(l.OutputHandler, logInfo, defaultOutputName, messages...)
}

func (l *Logger) Notice(messages ...any) {
	l.writeToLog(l.OutputHandler, logNotice, defaultOutputName, messages...)
}

func (l *Logger) Warn(messages ...any) {
	l.writeToLog(l.OutputHandler, logWarning, defaultOutputName, messages...)
}

func (l *Logger) WarnError(err error, messages ...any) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logWarning, defaultErrorName, messages...)
	strMessages := make([]string, len(messages))
	for i, msg := range messages {
		strMessages[i] = fmt.Sprintf("%v", msg)
	}
	return fmt.Errorf("%s", strings.Join(strMessages, " "))
}

func (l *Logger) Error(err error, messages ...any) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logError, defaultErrorName, messages...)
	strMessages := make([]string, len(messages))
	for i, msg := range messages {
		strMessages[i] = fmt.Sprintf("%v", msg)
	}
	return fmt.Errorf("%s", strings.Join(strMessages, " "))
}

func (l *Logger) Fatal(err error, messages ...any) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logFatal, defaultErrorName, messages...)
	strMessages := make([]string, len(messages))
	for i, msg := range messages {
		strMessages[i] = fmt.Sprintf("%v", msg)
	}
	return fmt.Errorf("%s", strings.Join(strMessages, " "))
}

func (l *Logger) Critical(err error, messages ...any) error {
	if err != nil {
		messages = append(messages, err.Error())
	}
	l.writeToLog(l.ErrorHandler, logCritical, defaultErrorName, messages...)
	strMessages := make([]string, len(messages))
	for i, msg := range messages {
		strMessages[i] = fmt.Sprintf("%v", msg)
	}
	return fmt.Errorf("%s", strings.Join(strMessages, " "))
}
