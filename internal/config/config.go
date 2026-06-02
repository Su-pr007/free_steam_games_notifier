package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	TgBot BotConfig
}

type AppConfig struct {
	Env string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      string
}

type BotConfig struct {
	Token     string
	ChannelId string
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	var cfg Config

	cfg.App.Env = getEnv("APP_ENV")
	cfg.DB.Host = getEnv("DB_HOST_FOR_PROJECT")
	cfg.DB.Port = getEnv("DB_PORT_FOR_PROJECT")
	cfg.DB.Username = getEnv("DB_USER")
	cfg.DB.Password = getEnv("DB_PASSWORD")
	cfg.DB.DBName = getEnv("DB_NAME")
	cfg.DB.SSL = getEnv("DB_SSL")
	cfg.TgBot.Token = getEnv("TG_BOT_TOKEN")

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load config from environment: " + err.Error())
	}

	return &cfg
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		panic("environment variable not found: " + key)
	}

	return value
}
