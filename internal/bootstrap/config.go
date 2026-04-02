package bootstrap

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.Logger, configPath, queryPath string, config any) {
	data, _ := os.ReadFile(configPath)
	expanded := os.ExpandEnv(string(data))

	logger.Debug("config", zap.String("expanded", expanded))

	viper.SetConfigType("yml")
	if err := viper.ReadConfig(strings.NewReader(expanded)); err != nil {
		logger.Fatal("error while reading config", zap.Error(err))
	}

	viper.SetConfigFile(queryPath)
	viper.MergeInConfig()

	if err := viper.Unmarshal(config); err != nil {
		logger.Fatal("error while reading config", zap.Error(err))
	}

	logger.Debug("config loaded", zap.Any("config", config))
}
