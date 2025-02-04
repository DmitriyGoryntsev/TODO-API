package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL         string
	ServerAddress string
}

func LoadConfig(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	return &Config{
		DBURL:         os.Getenv("DB_URL"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}, nil
}
