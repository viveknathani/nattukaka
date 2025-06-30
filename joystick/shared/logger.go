package shared

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger returns a new structured logger
func NewLogger() *zap.Logger {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevel(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "ts",
			EncodeTime:  zapcore.EpochMillisTimeEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		// not being able to setup the logger is a critical failure, exit the process in this case.
		log.Println(err)
		os.Exit(1)
	}

	return logger
}
