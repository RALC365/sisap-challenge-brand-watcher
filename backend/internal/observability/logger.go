package observability

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

func InitLogger() *zap.Logger {
	loggerOnce.Do(func() {
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.MillisDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		var err error
		logger, err = config.Build()
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})

	return logger
}

func GetLogger() *zap.Logger {
	if logger == nil {
		return InitLogger()
	}
	return logger
}

func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}
