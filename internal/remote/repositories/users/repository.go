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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Insert(tableName).
		Columns(columns[1:]...).
		Values(user.Name, pq.Array(user.Roles)).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	err = query.QueryRowContext(ctx).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, user users.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Update(tableName).
		Set("name", user.Name).
		Set("roles", pq.Array(user.Roles)).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	_, err = query.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *repository) Get(ctx context.Context, userID int64) (*users.User, error) {
	query := squirrel.Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	row := query.QueryRowContext(ctx)
	user := &users.User{}
	err := row.Scan(&user.ID, &user.Name, pq.Array(&user.Roles))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) List(ctx context.Context) ([]users.User, error) {
	query := squirrel.Select(columns...).
		From(tableName).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
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

	return usersList, nil
}

func (r *repository) Delete(ctx context.Context, userID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	_, err = query.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
