package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

type Logger = *zap.SugaredLogger

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

func New(level Level) Logger {
	core := zap.NewProductionConfig()
	core.Level = zap.NewAtomicLevelAt(getLevel(level))
	core.EncoderConfig.TimeKey = "timestamp"
	core.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := core.Build()
	return logger.Sugar()
}

func getLevel(level Level) zapcore.Level {
	switch level {
	case InfoLevel:
		return zapcore.InfoLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
