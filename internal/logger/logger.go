package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configuring the Logger used in Terrabutler
var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
}
var configZap = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	Development:      true,
	Encoding:         "console",
	EncoderConfig:    encoderConfig,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

var Zap, _ = configZap.Build()
