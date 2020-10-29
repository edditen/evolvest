package log

import (
	logging "log"
	"os"
)

var (
	defaultLog Log
)

func init() {
	defaultLog = NewConsole(false)
}

func SetVerbose(verbose bool) {
	defaultLog.SetVerbose(verbose)
}

type Verbose interface {
	SetVerbose(verbose bool)
}

type Log interface {
	Verbose
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Fatal(msg string, v ...interface{})
}

type Console struct {
	logger  *logging.Logger
	verbose bool
}

func NewConsole(verbose bool) *Console {
	return &Console{
		logger:  logging.New(os.Stdout, "", logging.LstdFlags),
		verbose: verbose,
	}
}

func (l *Console) SetVerbose(verbose bool) {
	l.verbose = verbose
}

func (l *Console) Debug(msg string, v ...interface{}) {
	if l.verbose {
		l.logger.Printf("[DEBUG] "+msg+"\n", v...)
	}
}

func (l *Console) Info(msg string, v ...interface{}) {
	l.logger.Printf("[INFO]  "+msg+"\n", v...)
}

func (l *Console) Warn(msg string, v ...interface{}) {
	l.logger.Printf("[WARN]  "+msg+"\n", v...)
}

func (l *Console) Fatal(msg string, v ...interface{}) {
	l.logger.Printf("[FATAL] "+msg+"\n", v...)
	os.Exit(1)
}

func Debug(msg string, v ...interface{}) {
	defaultLog.Debug(msg, v...)
}

func Info(msg string, v ...interface{}) {
	defaultLog.Info(msg, v...)
}

func Warn(msg string, v ...interface{}) {
	defaultLog.Warn(msg, v...)
}

func Fatal(msg string, v ...interface{}) {
	defaultLog.Fatal(msg, v...)
}
