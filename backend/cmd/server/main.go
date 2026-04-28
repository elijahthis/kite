package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elijahthis/kite/internal/config"
	"github.com/elijahthis/kite/internal/scripts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	f := config.NewFactory()

	if err := scripts.RunMigrations(f.Config.DB_MIGRATIONS_DIRECTORY, f.Config.DB_CONN_STR); err != nil {
		log.Error().Err(err).Msg("Unable to run migrations")
	}

	if err := scripts.SeedAdminUser(context.Background(), f.DB, f.Config.SYSTEM_USER_EMAIL, f.Config.SYSTEM_USER_PASSWORD, f.Config.SYSTEM_USER_NAME); err != nil {
		log.Error().Err(err).Msg("Unable to seed DB")
	}

	srv := &http.Server{
		Addr:         ":" + f.Config.PORT,
		Handler:      f.Router.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msgf("Server is listening on port %s\n", f.Config.PORT)
		if err := http.ListenAndServe(":"+f.Config.PORT, f.Router.Handler); err != nil {
			log.Fatal().Msg("Unable to start server")
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
		os.Exit(1)
	}

	log.Info().Msg("server exited cleanly")

}
