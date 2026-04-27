package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DB_CONN_STR             string
	DB_MIGRATIONS_DIRECTORY string
	JWT_SECRET_KEY          string
	PORT                    string
	SYSTEM_USER_NAME        string
	SYSTEM_USER_EMAIL       string
	SYSTEM_USER_PASSWORD    string
}

func mustEnv(v string) string {
	v, ok := os.LookupEnv(v)
	if !ok {
		log.Fatal().Msgf("%v is missing from env", v)

	}
	return v
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, relying on system environment variables")
	}

	PORT, exists := os.LookupEnv("PORT")
	if !exists {
		PORT = "8080"
	}

	return &Config{
		DB_CONN_STR:             mustEnv("DB_CONN_STR"),
		DB_MIGRATIONS_DIRECTORY: mustEnv("DB_MIGRATIONS_DIRECTORY"),
		JWT_SECRET_KEY:          mustEnv("JWT_SECRET_KEY"),
		PORT:                    PORT,
		SYSTEM_USER_NAME:        mustEnv("SYSTEM_USER_NAME"),
		SYSTEM_USER_EMAIL:       mustEnv("SYSTEM_USER_EMAIL"),
		SYSTEM_USER_PASSWORD:    mustEnv("SYSTEM_USER_PASSWORD"),
	}
}
