package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger type
type (
	Logger   = zap.SugaredLogger
	LogLevel = string
)

const (
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
	ErrorLevel LogLevel = "error"
)

// Supported log levels
var logLevels = map[string]zapcore.Level{
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	DebugLevel: zap.DebugLevel,
	ErrorLevel: zap.ErrorLevel,
}

var instance *zap.SugaredLogger
var once sync.Once

// Get the logger (Singleton pattern)
func getLogger() *Logger {
	once.Do(func() {
		// Do nothing if the logger was already instantiated.
		if instance == nil {
			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("Can't initialize zap logger: %v", err)
			}
			instance = logger.Sugar()
		}
	})
	return instance
}

// Error Log
func Error(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
}

// Info Log
func Info(template string, args ...interface{}) {
	getLogger().Infof(template, args...)
}

// Debug Log
func Debug(template string, args ...interface{}) {
	getLogger().Debugf(template, args...)
}

// Warn Log
func Warn(template string, args ...interface{}) {
	getLogger().Warnf(template, args...)
}

// Fatal Log
func Fatal(template string, args ...interface{}) {
	getLogger().Fatalf(template, args...)
}

// SetupLogger - Configure logger
func SetupLogger(location string, level LogLevel) {
	// Log level fallback
	if _, ok := logLevels[level]; !ok {
		level = InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   location,
			MaxBackups: 3, // 3 files
			MaxAge:     5, // 5 days
		}),
		logLevels[level],
	))
	instance = logger.Sugar()
}
