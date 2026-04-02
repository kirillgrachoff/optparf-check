package bootstrap

import "go.uber.org/zap"

func InitLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.DebugLevel)
	config.Encoding = "console"

	logger, err := config.Build()

	if err != nil {
		panic(err)
	}

	return logger
}