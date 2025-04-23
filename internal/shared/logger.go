package shared

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(appEnv string) (*zap.SugaredLogger, error) {
	var config zap.Config

	if appEnv == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		config.Encoding = "json"
	} else {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
