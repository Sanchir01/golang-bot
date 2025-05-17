package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

type Config struct {
	Env           string       `yaml:"env"`
	KafkaConsumer KafkaBrokers `yaml:"kafka"`
}

type KafkaBrokers struct {
	Consumer Consumer `yaml:"consumer"`
}
type Consumer struct {
	Topic   string   `yaml:"topic"`
	Broker  []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
}

func InitConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("ошибка при инициализации переменных окружения", err.Error())
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
