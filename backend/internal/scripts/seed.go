package scripts

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser(ctx context.Context, db *sqlx.DB, email, password, first_name string) error {
	var id uuid.UUID
	query := `
		INSERT INTO users(email, password_hash, first_name) 
		VALUES($1, $2, $3) 
		ON CONFLICT DO NOTHING
		RETURNING id;
	`

	h_bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password_hash := string(h_bytes)

	if err := db.QueryRowContext(ctx, query, email, password_hash, first_name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msg("Seed already present")
			return nil
		}
		return err
	}
	log.Info().Msg("DB seeded successfully!")

	return nil
}
