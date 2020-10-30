package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	debugEnable bool
)

func SetVerbose(verbose bool) {
	debugEnable = verbose
}

type Verbose interface {
	SetVerbose(verbose bool)
}

type Logger interface {
	Verbose
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Fatal(msg string, v ...interface{})
	WithField(field string, val interface{}) Logger
	WithError(err error) Logger
}

type Console struct {
	logger  *log.Logger
	verbose bool
	fields  map[string]interface{}
	err     error
}

func NewConsole(verbose bool) *Console {
	return &Console{
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		verbose: verbose,
		fields:  make(map[string]interface{}),
	}
}

func (l *Console) SetVerbose(verbose bool) {
	l.verbose = verbose
}

func (l *Console) Debug(msg string, v ...interface{}) {
	if l.verbose {
		l.logger.Println(l.build("DEBUG", msg, v...))
	}
}

func (l *Console) Info(msg string, v ...interface{}) {
	l.logger.Println(l.build("INFO", msg, v...))
}

func (l *Console) Warn(msg string, v ...interface{}) {
	l.logger.Println(l.build("WARN", msg, v...))
}

func (l *Console) Fatal(msg string, v ...interface{}) {
	l.logger.Println(l.build("FATAL", msg, v...))
	os.Exit(1)
}

func (l *Console) WithField(field string, val interface{}) Logger {
	l.checkNil()
	l.fields[field] = val
	return l
}

func (l *Console) WithError(err error) Logger {
	l.checkNil()
	l.err = err
	return l
}

func (l *Console) build(level, msg string, v ...interface{}) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] ", level))
	sb.WriteString(fmt.Sprintf(msg, v...))

	if len(l.fields) > 0 {
		sb.WriteString(" | ")
		for field, val := range l.fields {
			sb.WriteString(field)
			sb.WriteString("=")
			sb.WriteString(fmt.Sprintf("%v", val))
			sb.WriteString(",")
		}
	}

	if l.err != nil {
		sb.WriteString(" | ")
		sb.WriteString("err")
		sb.WriteString("=")
		sb.WriteString(l.err.Error())
	}
	return sb.String()
}

func (l *Console) checkNil() {
	if l.fields == nil {
		l.fields = make(map[string]interface{})
	}
}

func Debug(msg string, v ...interface{}) {
	NewConsole(debugEnable).Debug(msg, v...)
}

func Info(msg string, v ...interface{}) {
	NewConsole(debugEnable).Info(msg, v...)
}

func Warn(msg string, v ...interface{}) {
	NewConsole(debugEnable).Warn(msg, v...)
}

func Fatal(msg string, v ...interface{}) {
	NewConsole(debugEnable).Fatal(msg, v...)
}

func WithField(field string, val interface{}) Logger {
	l := NewConsole(debugEnable)
	l.WithField(field, val)
	return l
}

func WithError(err error) Logger {
	l := NewConsole(debugEnable)
	l.WithError(err)
	return l
}
