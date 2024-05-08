package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRoles, downCreateRoles)
}

func upCreateRoles(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT NULL,
			read BOOLEAN NOT NULL DEFAULT FALSE,
			write BOOLEAN NOT NULL DEFAULT FALSE,
			delete BOOLEAN NOT NULL DEFAULT FALSE,
			overwrite BOOLEAN NOT NULL DEFAULT FALSE,
			superuser BOOLEAN NOT NULL DEFAULT FALSE,
			spaces TEXT[] NULL
		);
		CREATE INDEX IF NOT EXISTS idx_roles_id ON roles (id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles (name);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateRoles(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS roles;
		DROP INDEX IF EXISTS idx_roles_id;
		DROP INDEX IF EXISTS idx_roles_name;
	`)
	if err != nil {
		return err
	}
	return nil
}
