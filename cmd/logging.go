package cmd

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string) (*zap.Logger, error) {
	zapCfg := zap.NewProductionEncoderConfig()
	zapCfg.TimeKey = "time"
	zapCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	pCfg := zap.NewProductionConfig()
	logLevel, err := logLevelFromCfg(level)
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}
	pCfg.Level = zap.NewAtomicLevelAt(logLevel)
	pCfg.EncoderConfig = zapCfg
	pCfg.DisableStacktrace = true

	logger, err := pCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}

	return logger, nil
}

func logLevelFromCfg(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "info":
		fallthrough
	case "":
		return zap.InfoLevel, nil
	default:
		return zap.InfoLevel, errors.New("unknown log level: " + level)
	}
}
