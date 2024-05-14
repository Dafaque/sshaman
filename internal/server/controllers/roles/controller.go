package roles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dafaque/sshaman/internal/server/errs"
	"github.com/lib/pq"
)

const SUID int64 = 1

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
	isPermitted, err := rc.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	return rc.rolesRepository.Create(ctx, role)
}

func (rc *controller) Update(ctx context.Context, role *Role) error {
	isPermitted, err := rc.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	if role.ID == SUID {
		return errors.New("can't update superuser role")
	}
	err = rc.rolesRepository.Update(ctx, role)
	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code == "23505" {
			return fmt.Errorf("role %s already exists", role.Name)
		}
	}
	return err
}

func (rc *controller) Delete(ctx context.Context, id int64) error {
	isPermitted, err := rc.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	if id == SUID {
		return errors.New("can't delete superuser role")
	}
	return rc.rolesRepository.Delete(ctx, id)
}

func (rc *controller) List(ctx context.Context) ([]Role, error) {
	isPermitted, err := rc.isPermitted(ctx)
	if err != nil {
		return nil, err
	}
	if !isPermitted {
		return nil, errs.ErrNotPermitted
	}
	return rc.rolesRepository.List(ctx)
}

func (rc *controller) Get(ctx context.Context, ids ...int64) ([]Role, error) {
	//? not exposed by the API, internal usage
	return rc.rolesRepository.Get(ctx, ids...)
}

var errUnshureSuperuserRole = errors.New("failed to ensure superuser role")

func (rc *controller) enshureSuperuser() error {
	ctx := context.WithValue(
		context.Background(),
		"permissions",
		internalOperationsPermissions{},
	)
	su, err := rc.Get(ctx, SUID)
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
	if err := rc.Create(ctx, role); err != nil {
		return errors.Join(errUnshureSuperuserRole, err)
	}
	return nil
}

func (rc *controller) isPermitted(ctx context.Context) (bool, error) {
	perms, ok := ctx.Value("permissions").(permissions)
	if !ok {
		return false, errors.New("permissions not found")
	}
	return perms.SU(), nil
}
