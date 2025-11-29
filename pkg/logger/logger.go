package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(development bool) error {
	var config zap.Config
	if development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var err error
	globalLogger, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(globalLogger)
	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		// Fallback to development logger if not initialized
		Init(true)
	}
	return globalLogger
}

func Sync() {
	if globalLogger != nil {
		globalLogger.Sync()
	}
}
