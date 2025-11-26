package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBName           string `mapstructure:"DB_NAME"`
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPort        string `mapstructure:"REDIS_PORT"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
	ServerPort       string `mapstructure:"SERVER_PORT"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	JWTRefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
	JWTAccessExpiry  string `mapstructure:"JWT_ACCESS_EXPIRY"`
	JWTRefreshExpiry string `mapstructure:"JWT_REFRESH_EXPIRY"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we might be using env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	err := viper.Unmarshal(&config)
	return config, err
}
