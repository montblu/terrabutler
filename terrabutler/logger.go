package main

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
var config = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:      true,
	Encoding:         "console",
	EncoderConfig:    encoderConfig,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}
var logger, _ = config.Build()
