package logging

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitializeZapLogger(level string) error {
	if Log != nil {
		return nil
	}
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl.Sugar()
	return nil
}
