package main

import (
	// "github.com/elijahthis/kite/internal/database"
	// "github.com/elijahthis/kite/internal/factory"

	"os"
	"time"

	"github.com/elijahthis/kite/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	f := config.NewFactory()

	// apps
	// handlers

	// if err := database.RunMigrations(f.DBConnStr); err != nil {
	// 	log.Error().Err(err).Msg("Unable to run migrations")
	// }
}
