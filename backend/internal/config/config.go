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
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, relying on system environment variables")
	}

	DB_CONN_STR, exists := os.LookupEnv("DB_CONN_STR")
	if !exists {
		log.Fatal().Msg("DB_CONN_STR is missing from env")
	}

	DB_MIGRATIONS_DIRECTORY, exists := os.LookupEnv("DB_MIGRATIONS_DIRECTORY")
	if !exists {
		log.Fatal().Msg("DB_MIGRATIONS_DIRECTORY is missing from env")
	}

	JWT_SECRET_KEY, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		log.Fatal().Msg("JWT_SECRET_KEY is missing from env")
	}

	PORT, exists := os.LookupEnv("PORT")
	if !exists {
		PORT = "8080"
	}

	return &Config{
		DB_CONN_STR:             DB_CONN_STR,
		DB_MIGRATIONS_DIRECTORY: DB_MIGRATIONS_DIRECTORY,
		JWT_SECRET_KEY:          JWT_SECRET_KEY,
		PORT:                    PORT,
	}
}
