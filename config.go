package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	Port     uint   `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}

	if cfg.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return &cfg, nil
}
