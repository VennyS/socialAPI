package shared

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger инициализирует логгер в зависимости от окружения
func InitLogger(appEnv string) (*zap.SugaredLogger, error) {
	var config zap.Config

	// Настроим логгер для production и development
	if appEnv == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		config.Encoding = "json"
	} else {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
