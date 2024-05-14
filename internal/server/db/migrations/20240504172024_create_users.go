package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsers, downCreateUsers)
}

func upCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			roles INTEGER[] NOT NULL
		);
		COMMENT ON TABLE users IS 'Users';
		COMMENT ON COLUMN users.id IS 'The unique identifier for the user';
		COMMENT ON COLUMN users.name IS 'The name for the user';
		COMMENT ON COLUMN users.roles IS 'The roles for the user';
		
		CREATE INDEX IF NOT EXISTS idx_users_id ON users (id);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS users;
		DROP INDEX IF EXISTS idx_users_id;
	`)
	if err != nil {
		return err
	}
	return nil
}
