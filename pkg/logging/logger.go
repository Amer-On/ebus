package logging

import "go.uber.org/zap"

// TODO: Provide usefull settings to configure logger
func InitLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	return logger, err
}
