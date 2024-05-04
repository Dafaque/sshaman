package users

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/Dafaque/sshaman/internal/remote/controllers/users"
)

const tableName = "users"

var columns = []string{
	"id",
	"name",
	"roles",
}

type Repository interface {
	Create(ctx context.Context, user *users.User) error
	Update(ctx context.Context, user users.User) error
	Get(ctx context.Context, userID int64) (*users.User, error)
	List(ctx context.Context) ([]users.User, error)
	Delete(ctx context.Context, userID int64) error
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *users.User) error {
	query := squirrel.Insert(tableName).
		Columns(columns...).
		Values(user.ID, user.Name, pq.Array(user.Roles)).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		RunWith(r.db).
		QueryRowContext(ctx)

	user.ID = -1
	err := query.Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, user users.User) error {
	result, err := squirrel.Update(tableName).
		Set("name", user.Name).
		Set("roles", pq.Array(user.Roles)).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db).
		ExecContext(ctx)

	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, userID int64) (*users.User, error) {
	query := squirrel.Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db).
		QueryRowContext(ctx)

	var user users.User
	err := query.Scan(&user.ID, &user.Name, pq.Array(&user.Roles))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) List(ctx context.Context) ([]users.User, error) {
	rows, err := squirrel.Select(columns...).
		From(tableName).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersList []users.User
	for rows.Next() {
		var user users.User
		err := rows.Scan(&user.ID, &user.Name, pq.Array(&user.Roles))
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usersList, nil
}

func (r *repository) Delete(ctx context.Context, userID int64) error {
	result, err := squirrel.Delete(tableName).
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db).
		ExecContext(ctx)

	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
