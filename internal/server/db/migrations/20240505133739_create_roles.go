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

		COMMENT ON TABLE roles IS 'Roles for users';
		COMMENT ON COLUMN roles.id IS 'The unique identifier for the role';
		COMMENT ON COLUMN roles.name IS 'The name for the role; Must be unique';
		COMMENT ON COLUMN roles.description IS 'The description for the role';
		COMMENT ON COLUMN roles.read IS 'The read permission for the role';
		COMMENT ON COLUMN roles.write IS 'The write permission for the role';
		COMMENT ON COLUMN roles.delete IS 'The delete permission for the role';
		COMMENT ON COLUMN roles.overwrite IS 'The overwrite permission for the role';
		COMMENT ON COLUMN roles.superuser IS 'The superuser permission for the role';
		COMMENT ON COLUMN roles.spaces IS 'The spaces for the role';

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
