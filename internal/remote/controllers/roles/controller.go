package roles

import (
	"context"
	"database/sql"
	"errors"
)

type Controller interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]Role, error)
	Get(ctx context.Context, ids ...int64) ([]Role, error)
}

type controller struct {
	rolesRepository repository
}

func New(rolesRepository repository) (Controller, error) {
	cnt := &controller{rolesRepository: rolesRepository}
	if err := cnt.enshureSuperuser(); err != nil {
		return nil, err
	}
	return cnt, nil
}

func (rc *controller) Create(ctx context.Context, role *Role) error {
	return rc.rolesRepository.Create(ctx, role)
}

func (rc *controller) Update(ctx context.Context, role *Role) error {
	return rc.rolesRepository.Update(ctx, role)
}

func (rc *controller) Delete(ctx context.Context, id int64) error {
	return rc.rolesRepository.Delete(ctx, id)
}

func (rc *controller) List(ctx context.Context) ([]Role, error) {
	return rc.rolesRepository.List(ctx)
}

func (rc *controller) Get(ctx context.Context, ids ...int64) ([]Role, error) {
	return rc.rolesRepository.Get(ctx, ids...)
}

var errUnshureSuperuserRole = errors.New("failed to ensure superuser role")

func (rc *controller) enshureSuperuser() error {
	su, err := rc.Get(context.Background(), 0)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Join(errUnshureSuperuserRole, err)
		}
	}
	if len(su) > 0 {
		return nil
	}
	role := &Role{
		Name:        "su",
		Description: "Full permissions role",
		Read:        true,
		Write:       true,
		Delete:      true,
		Overwrite:   true,
		SU:          true,
		Spaces:      []string{"*"},
	}
	if err := rc.Create(context.Background(), role); err != nil {
		return errors.Join(errUnshureSuperuserRole, err)
	}
	return nil
}
