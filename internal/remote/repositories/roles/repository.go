package roles

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/Dafaque/sshaman/internal/remote/controllers/roles"
)

const tableName = "roles"

var columns = []string{
	"id",
	"name",
	"description",
	"read",
	"write",
	"delete",
	"overwrite",
	"superuser",
	"spaces",
}

type Repository interface {
	Create(ctx context.Context, role *roles.Role) error
	Update(ctx context.Context, role *roles.Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]roles.Role, error)
	Get(ctx context.Context, ids ...int64) ([]roles.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *roles.Role) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Insert(tableName).
		Columns(columns[1:]...).
		Values(role.Name,
			role.Description,
			role.Read,
			role.Write,
			role.Delete,
			role.Overwrite,
			role.SU,
			pq.Array(role.Spaces),
		).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	err = query.QueryRowContext(ctx).Scan(&role.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *roleRepository) Update(ctx context.Context, role *roles.Role) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Update(tableName).
		Set("name", role.Name).
		Set("description", role.Description).
		Set("read", role.Read).
		Set("write", role.Write).
		Set("delete", role.Delete).
		Set("overwrite", role.Overwrite).
		Set("su", role.SU).
		Set("spaces", pq.Array(role.Spaces)).
		Where(squirrel.Eq{"id": role.ID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	_, err = query.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := squirrel.Delete(tableName).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	_, err = query.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *roleRepository) List(ctx context.Context) ([]roles.Role, error) {
	query := squirrel.Select(columns...).
		From(tableName).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []roles.Role
	for rows.Next() {
		var role roles.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Read, &role.Write, &role.Delete, &role.Overwrite, &role.SU, pq.Array(&role.Spaces))
		if err != nil {
			return nil, err
		}
		list = append(list, role)
	}
	return list, nil
}

func (r *roleRepository) Get(ctx context.Context, ids ...int64) ([]roles.Role, error) {
	query := squirrel.Select(columns...).
		From(tableName).
		Where(squirrel.Eq{"id": ids}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolesList []roles.Role
	for rows.Next() {
		var role roles.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Read, &role.Write, &role.Delete, &role.Overwrite, &role.SU, pq.Array(&role.Spaces))
		if err != nil {
			return nil, err
		}
		rolesList = append(rolesList, role)
	}

	return rolesList, nil
}
