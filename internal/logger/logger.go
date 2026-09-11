package logger

import "go.uber.org/zap"

const LogLevel = "info"

func NewLogger() (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(LogLevel)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl, nil
}
