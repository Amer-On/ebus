package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	// Убираем дату, таймзону — оставляем только точное локальное время с миллисекундами
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	// cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	logger, err := cfg.Build()
	return logger, err
}
