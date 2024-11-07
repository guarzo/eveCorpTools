package xlog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// XLogger is a custom logger struct that wraps a logrus.Logger and retains xlog-specific features.
type XLogger struct {
	log *logrus.Logger
}

// NewXLogger initializes a new XLogger with the provided logrus.Logger instance.
func NewXLogger(log *logrus.Logger) *XLogger {
	return &XLogger{log: log}
}

func getFilePathWithDir(file string) string {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return file
	}
	relativePath, err := filepath.Rel(projectRoot, file)
	if err != nil {
		fmt.Println("Error getting relative path:", err)
		return file
	}
	return relativePath
}

// logMessage is a helper function that logs a message with its location
func logMessage(message string, file string, line int, ok bool) {
	if ok {
		fmt.Printf("%s:%d - %s\n", getFilePathWithDir(file), line, message)
	} else {
		fmt.Printf("%s\n", message)
	}
}

func Log(message string) {
	_, file, line, ok := runtime.Caller(1)
	logMessage(message, file, line, ok)
}

func Logf(format string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	logMessage(fmt.Sprintf(format, a...), file, line, ok)
}

func LogIndirect(message string) {
	_, file, line, ok := runtime.Caller(2)
	logMessage(message, file, line, ok)
}

func Logt(functionName string, start time.Time) {
	_, file, line, ok := runtime.Caller(1)
	logMessage(fmt.Sprintf("%s - %vms", functionName, time.Since(start).Milliseconds()), file, line, ok)
}

//package xlog
//
//import (
//"fmt"
//"os"
//"path/filepath"
//"runtime"
//"time"
//
//"github.com/sirupsen/logrus"
//)
//
//type XLogger struct {
//	log *logrus.Logger
//}
//
//// NewXLogger initializes a new XLogger with the provided logrus.Logger instance.
//func NewXLogger(logger *logrus.Logger) *XLogger {
//	return &XLogger{log: logger}
//}
//
//func getFilePathWithDir(file string) string {
//	projectRoot, err := os.Getwd()
//	if err != nil {
//		fmt.Println("Error getting current directory:", err)
//		return file
//	}
//	relativePath, err := filepath.Rel(projectRoot, file)
//	if err != nil {
//		fmt.Println("Error getting relative path:", err)
//		return file
//	}
//	return relativePath
//}
//
//// logMessage is a helper function that logs a message with its location using logrus.
//func (l *XLogger) logMessage(message string, file string, line int, ok bool) {
//	entry := l.log.WithFields(logrus.Fields{
//		"file": getFilePathWithDir(file),
//		"line": line,
//	})
//	if ok {
//		entry.Info(message)
//	} else {
//		l.log.Info(message)
//	}
//}
//
//// Log logs a message with file path and line number.
//func (l *XLogger) Log(message string) {
//	_, file, line, ok := runtime.Caller(1)
//	l.logMessage(message, file, line, ok)
//}
//
//// Logf logs a formatted message with file path and line number.
//func (l *XLogger) Logf(format string, a ...interface{}) {
//	_, file, line, ok := runtime.Caller(1)
//	l.logMessage(fmt.Sprintf(format, a...), file, line, ok)
//}
//
//// LogIndirect logs a message with file path and line number, adjusted for indirect logging.
//func (l *XLogger) LogIndirect(message string) {
//	_, file, line, ok := runtime.Caller(2)
//	l.logMessage(message, file, line, ok)
//}
//
//// Logt logs a message with the function name and execution time in milliseconds.
//func (l *XLogger) Logt(functionName string, start time.Time) {
//	_, file, line, ok := runtime.Caller(1)
//	elapsed := time.Since(start).Milliseconds()
//	l.logMessage(fmt.Sprintf("%s - %vms", functionName, elapsed), file, line, ok)
//}
