package db

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/Dafaque/sshaman/internal/remote/config"
)

func New(cfg *config.Config) (*sql.DB, error) {
	var connstr bytes.Buffer
	connstr.WriteString("host=")
	connstr.WriteString(cfg.Postgres.Host)
	connstr.WriteString(" port=")
	connstr.WriteString(fmt.Sprintf("%d", cfg.Postgres.Port))
	connstr.WriteString(" user=")
	connstr.WriteString(cfg.Postgres.User)
	connstr.WriteString(" password=")
	connstr.WriteString(cfg.Postgres.Password)
	connstr.WriteString(" dbname=")
	connstr.WriteString(cfg.Postgres.DBName)
	connstr.WriteString(" sslmode=")
	connstr.WriteString(cfg.Postgres.SSLMode)

	db, err := sql.Open("postgres", connstr.String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
