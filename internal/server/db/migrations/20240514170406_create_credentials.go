package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCredentials, downCreateCredentials)
}

func upCreateCredentials(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE credentials (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id) ON DELETE NULL,
			space TEXT NOT NULL,
			alias TEXT NOT NULL,
			username TEXT NOT NULL,
			host TEXT NOT NULL,
			port INTEGER NOT NULL,
			password TEXT,
			key BYTEA,
			passphrase BYTEA,
		);
		COMMENT ON TABLE credentials IS 'Credentials for SSH connections';
		COMMENT ON COLUMN credentials.id IS 'The unique identifier for the credential';
		COMMENT ON COLUMN credentials.user_id IS 'The user identifier for the credential';
		COMMENT ON COLUMN credentials.space IS 'The space for the credential';
		COMMENT ON COLUMN credentials.alias IS 'The alias for the credential';
		COMMENT ON COLUMN credentials.username IS 'The username for the credential';
		COMMENT ON COLUMN credentials.host IS 'The host for the credential';
		COMMENT ON COLUMN credentials.port IS 'The port for the credential';
		COMMENT ON COLUMN credentials.password IS 'The password for the credential';
		COMMENT ON COLUMN credentials.key IS 'The key for the credential';
		COMMENT ON COLUMN credentials.passphrase IS 'The passphrase for the credential';

		CREATE INDEX idx_credentials_user_id ON credentials (user_id);
		CREATE INDEX idx_credentials_alias ON credentials (alias);
		CREATE INDEX idx_credentials_space ON credentials (space);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateCredentials(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS credentials;
		DROP INDEX IF EXISTS idx_credentials_user_id;
		DROP INDEX IF EXISTS idx_credentials_alias;
		DROP INDEX IF EXISTS idx_credentials_space;
	`)
	if err != nil {
		return err
	}
	return nil
}
