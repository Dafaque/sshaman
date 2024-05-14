package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dafaque/sshaman/internal/server/errs"
)

// supersuser id
const SUID int64 = 1

type Controller interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user User) error
	Get(ctx context.Context, userID int64) (*User, error)
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context) ([]User, error)
	PrintToken() bool
}

type controller struct {
	usersRepository usersRepository
	suCreated       bool
}

func New(usersRepository usersRepository) (Controller, error) {
	controller := &controller{
		usersRepository: usersRepository,
	}
	err := controller.enshureSuperuser()
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func (c *controller) Create(ctx context.Context, user *User) error {
	isPermitted, err := c.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	user.ID = -1
	return c.usersRepository.Create(ctx, user)
}

func (c *controller) Update(ctx context.Context, user User) error {
	isPermitted, err := c.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	if user.ID == SUID {
		return errors.New("can't update superuser")
	}
	return c.usersRepository.Update(ctx, user)
}

func (c *controller) Get(ctx context.Context, userID int64) (*User, error) {
	//? not exposed by the API, internal usage
	return c.usersRepository.Get(ctx, userID)
}

func (c *controller) Delete(ctx context.Context, userID int64) error {
	isPermitted, err := c.isPermitted(ctx)
	if err != nil {
		return err
	}
	if !isPermitted {
		return errs.ErrNotPermitted
	}
	if userID == SUID {
		return errors.New("can't delete superuser")
	}
	return c.usersRepository.Delete(ctx, userID)
}

func (c *controller) List(ctx context.Context) ([]User, error) {
	isPermitted, err := c.isPermitted(ctx)
	if err != nil {
		return nil, err
	}
	if !isPermitted {
		return nil, errs.ErrNotPermitted
	}
	return c.usersRepository.List(ctx)
}

var errEnshureSuperuser = errors.New("failed to ensure superuser")

func (c *controller) enshureSuperuser() error {
	ctx := context.WithValue(
		context.Background(),
		"permissions",
		internalOperationsPermissions{},
	)
	su, err := c.Get(ctx, SUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Join(errEnshureSuperuser, err)
		}
	}
	if su != nil {
		return nil
	}
	user := &User{
		Name:  "superuser",
		Roles: []int64{SUID},
	}
	if err := c.Create(ctx, user); err != nil {
		return errors.Join(errEnshureSuperuser, err)
	}
	c.suCreated = true
	return nil
}

func (c *controller) PrintToken() bool {
	return c.suCreated
}

func (rc *controller) isPermitted(ctx context.Context) (bool, error) {
	perms, ok := ctx.Value("permissions").(permissions)
	if !ok {
		return false, errors.New("permissions not found")
	}
	return perms.SU(), nil
}
