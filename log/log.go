package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	DEFAULT_CALL_DEPTH        = 3
	DEFAULT_LOGGER_CALL_DEPTH = 3
)

//var DefaultLogFlag = log.Lmicroseconds | log.Llongfile
var DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(LEVEL_INFO, os.Stdout, DefaultLogFlag)
}

func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

func NewLogger(level Level, writer io.Writer, flag int) *Logger {

	baseLogger := log.New(writer, "", flag)

	logger := &Logger{level: level, flag: flag}

	logger.AddLogger(baseLogger)

	return logger
}

type Logger struct {
	prefix string

	writerList []*log.Logger

	level Level

	buf bytes.Buffer

	flag int
}

func (l *Logger) AddLogger(logger *log.Logger) {

	l.writerList = append(l.writerList, logger)
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix

	for _, writer := range l.writerList {
		writer.SetPrefix(prefix)
	}
}

func (l *Logger) SetFlags(flag int) {

	for _, writer := range l.writerList {
		writer.SetFlags(flag)
	}
}

// 实现io.Writer接口
func (l *Logger) Write(level Level, msg string, callerSkip int) {

	if len(l.writerList) < 1 {
		return
	}

	l.buf.Reset()

	l.buf.WriteString(level.String())
	l.buf.WriteByte(' ')
	l.buf.WriteString(msg)

	for _, writer := range l.writerList {
		writer.Output(callerSkip, l.buf.String())
	}

	return
}

func (l *Logger) Log(level Level, msg string) {
	if l.level > level {
		return
	}

	l.Write(level, msg, DEFAULT_CALL_DEPTH)
}

func (l *Logger) Logf(level Level, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	l.Write(level, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
}

func (l *Logger) Debug(msg string) {
	if l.level > LEVEL_DEBUG {
		return
	}

	l.Write(LEVEL_DEBUG, msg, DEFAULT_CALL_DEPTH)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level > LEVEL_DEBUG {
		return
	}

	l.Write(LEVEL_DEBUG, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
}

func (l *Logger) Info(msg string) {
	if l.level > LEVEL_INFO {
		return
	}

	l.Write(LEVEL_INFO, msg, DEFAULT_CALL_DEPTH)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > LEVEL_INFO {
		return
	}

	l.Write(LEVEL_INFO, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
}

func (l *Logger) Warn(msg string) {
	if l.level > LEVEL_WARN {
		return
	}

	l.Write(LEVEL_WARN, msg, DEFAULT_CALL_DEPTH)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level > LEVEL_WARN {
		return
	}

	l.Write(LEVEL_WARN, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
}

func (l *Logger) Error(msg string) {
	l.Write(LEVEL_ERROR, msg, DEFAULT_CALL_DEPTH)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Write(LEVEL_ERROR, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
}

func (l *Logger) Fatal(msg string) {

	l.Write(LEVEL_FATAL, msg, DEFAULT_CALL_DEPTH)

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Write(LEVEL_FATAL, fmt.Sprintf(format, v...), DEFAULT_CALL_DEPTH)
	os.Exit(1)
}

/*****************************************************/

func Debug(msg string) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_DEBUG {
		return
	}

	defaultLogger.Write(LEVEL_DEBUG, msg, DEFAULT_LOGGER_CALL_DEPTH)
}

func Debugf(format string, v ...interface{}) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_DEBUG {
		return
	}

	defaultLogger.Write(LEVEL_DEBUG, fmt.Sprintf(format, v...), DEFAULT_LOGGER_CALL_DEPTH)
}

func Info(msg string) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_INFO {
		return
	}

	defaultLogger.Write(LEVEL_INFO, msg, DEFAULT_LOGGER_CALL_DEPTH)
}

func Infof(format string, v ...interface{}) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_INFO {
		return
	}

	defaultLogger.Write(LEVEL_INFO, fmt.Sprintf(format, v...), DEFAULT_LOGGER_CALL_DEPTH)
}

func Warn(msg string) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_WARN {
		return
	}

	defaultLogger.Write(LEVEL_WARN, msg, DEFAULT_LOGGER_CALL_DEPTH)
}

func Warnf(format string, v ...interface{}) {
	if defaultLogger == nil {
		return
	}
	if defaultLogger.level > LEVEL_WARN {
		return
	}

	defaultLogger.Write(LEVEL_WARN, fmt.Sprintf(format, v...), DEFAULT_LOGGER_CALL_DEPTH)
}

func Error(msg string) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Write(LEVEL_ERROR, msg, DEFAULT_LOGGER_CALL_DEPTH)
}

func Errorf(format string, v ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Write(LEVEL_ERROR, fmt.Sprintf(format, v...), DEFAULT_LOGGER_CALL_DEPTH)
}

func Fatal(msg string) {

	if defaultLogger != nil {
		defaultLogger.Write(LEVEL_FATAL, msg, DEFAULT_LOGGER_CALL_DEPTH)
	}

	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Write(LEVEL_FATAL, fmt.Sprintf(format, v...), DEFAULT_LOGGER_CALL_DEPTH)
	}
	os.Exit(1)
}
