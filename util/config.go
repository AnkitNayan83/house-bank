package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	GinServerAddress       string        `mapstructure:"GIN_SERVER_ADDRESS"`
	HttpServerAddress      string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress      string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TOKEN_SYMMETRIC_KEY    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Environment            string        `mapstructure:"ENVIRONMENT"`
}

func LoadConfig(path string) (config Config, err error) {
	// Read the config file
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
