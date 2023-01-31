package log

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var lvl = os.Getenv("LOG_LEVEL")

func New() *zap.SugaredLogger {
	if lvl == "" {
		lvl = "info"
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	level := zapLogger.Level()
	if err := level.Set(lvl); err != nil {
		_ = level.Set("info")
	}

	return zapLogger.Sugar()
}
