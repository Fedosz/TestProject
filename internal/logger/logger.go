package logger

import (
	"go.uber.org/zap"
)

func MustNew() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}
