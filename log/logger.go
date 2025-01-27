package log

import (
	"time"
)

type Factory interface {
	Get(name string, config interface{}) (*Logger, error)
}
type Builder interface {
	Logln(level Level, args ...interface{})
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	WithError(err error) Builder
	WithField(key string, value interface{}) Builder
	WithFields(fields Fields) Builder
}

type EntryTransporter interface {
	Transport(entry *Entry) error
	Level() Level
	Close() error
}

type Logger struct {
	transporter EntryTransporter
	// Flag for whether to log caller info (off by default)
	reportCaller bool
	packageName  string

	exitFunc exitFunc
}

func NewLogger(transporter EntryTransporter, reportCaller bool, packageName string) *Logger {
	return &Logger{transporter: transporter, reportCaller: reportCaller, packageName: packageName}
}

func (logger *Logger) SetTransporter(transporter EntryTransporter) {
	logger.transporter = transporter
}

func (logger *Logger) Logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.newBuilder().Logln(level, args...)
	}
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.newBuilder().Log(level, args...)
	}
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.newBuilder().Logf(level, format, args...)
	}
}

func (logger *Logger) Transport(entry *Entry) (err error) {
	if logger.IsLevelEnabled(entry.Level) {
		err = logger.transporter.Transport(entry)
	}
	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if entry.Level <= PanicLevel {
		panic(entry)
	}
	return err
}

func (logger *Logger) newBuilder() Builder {
	return &EntryBuilder{
		logger: logger,
		Time:   time.Now(),
	}
}
func (logger *Logger) WithError(err error) Builder {
	return logger.newBuilder().WithError(err)
}

func (logger *Logger) WithField(key string, value interface{}) Builder {
	return logger.newBuilder().WithField(key, value)
}

func (logger *Logger) WithFields(fields Fields) Builder {
	return logger.newBuilder().WithFields(fields)
}

func (logger *Logger) level() Level {
	return logger.transporter.Level()
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}
