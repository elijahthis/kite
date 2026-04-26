package scripts

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func RunMigrations(migrationsDir string, connStr string) error {
	fmt.Println(migrationsDir)
	log.Info().Msg("Running DB migrations...")

	m, err := migrate.New(migrationsDir, connStr)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("DB schema is already up to date.")
			return nil
		}
		return err
	}

	log.Info().Msg("DB migrations applied successfully!")

	return nil
}
