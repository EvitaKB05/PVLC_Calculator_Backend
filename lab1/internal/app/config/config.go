package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
}

func NewConfig() (*Config, error) {
	// Загружаем .env файл
	_ = godotenv.Load()

	// Настраиваем Viper для чтения config.toml
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	// Перезаписываем значения из переменных окружения
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.DBHost = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		cfg.DBPort = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.DBUser = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.DBPassword = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.DBName = name
	}

	log.Info("Конфигурация успешно загружена")
	return cfg, nil
}
