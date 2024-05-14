package credentials

import (
	"context"
	"database/sql"

	"github.com/Dafaque/sshaman/internal/server/controllers/credentials"
	"github.com/Masterminds/squirrel"
)

const (
	tableName = "credentials"
)

var fields = []string{
	"id",
	"user_id",
	"space",
	"alias",
	"username",
	"host",
	"port",
	"password",
	"key",
	"passphrase",
}

type Repository interface {
	Get(ctx context.Context, cred *credentials.Credentials) (*credentials.Credentials, error)
	Create(ctx context.Context, cred *credentials.Credentials) error
	Update(ctx context.Context, cred *credentials.Credentials) error
	List(ctx context.Context, userID int64) ([]*credentials.Credentials, error)
	Delete(ctx context.Context, cred *credentials.Credentials) error
	Drop(ctx context.Context, userID int64) error
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Get(ctx context.Context, cred *credentials.Credentials) (*credentials.Credentials, error) {
	query := squirrel.Select(fields...).
		From(tableName).
		Where(squirrel.Eq{
			"user_id": cred.UserID,
			"space":   cred.Space,
			"id":      cred.ID,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, sql, args...)
	var c credentials.Credentials
	err = row.Scan(
		&c.ID,
		&c.UserID,
		&c.Space,
		&c.Alias,
		&c.Username,
		&c.Host,
		&c.Port,
		&c.Password,
		&c.Key,
		&c.Passphrase,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) Create(ctx context.Context, cred *credentials.Credentials) error {
	query := squirrel.Insert(tableName).
		Columns(fields[1:]...). //? Skip ID as it is auto-generated
		Values(
			cred.UserID,
			cred.Space,
			cred.Alias,
			cred.Username,
			cred.Host,
			cred.Port,
			cred.Password,
			cred.Key,
			cred.Passphrase,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, sql, args...).Scan(&cred.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, cred *credentials.Credentials) error {
	setMap := map[string]interface{}{}

	if cred.Username != "" {
		setMap["username"] = cred.Username
	}
	if cred.Host != "" {
		setMap["host"] = cred.Host
	}
	if cred.Port != 0 {
		setMap["port"] = cred.Port
	}
	setMap["password"] = cred.Password
	setMap["key"] = cred.Key
	setMap["passphrase"] = cred.Passphrase

	query := squirrel.Update(tableName).
		SetMap(setMap).
		Where(squirrel.Eq{
			"user_id": cred.UserID,
			"space":   cred.Space,
			"alias":   cred.Alias,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) List(ctx context.Context, userID int64) ([]*credentials.Credentials, error) {
	query := squirrel.Select(fields...).
		From(tableName).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cs []*credentials.Credentials
	for rows.Next() {
		var c credentials.Credentials
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Space,
			&c.Alias,
			&c.Username,
			&c.Host,
			&c.Port,
			&c.Password,
			&c.Key,
			&c.Passphrase,
		); err != nil {
			return nil, err
		}
		cs = append(cs, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (r *repository) Delete(ctx context.Context, cred *credentials.Credentials) error {
	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{
			"id":      cred.ID,
			"user_id": cred.UserID,
			"space":   cred.Space,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) Drop(ctx context.Context, userID int64) error {
	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
