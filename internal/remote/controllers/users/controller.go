package users

import (
	"context"
	"database/sql"
	"errors"
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
	user.ID = -1
	return c.usersRepository.Create(ctx, user)
}

func (c *controller) Update(ctx context.Context, user User) error {
	return c.usersRepository.Update(ctx, user)
}

func (c *controller) Get(ctx context.Context, userID int64) (*User, error) {
	return c.usersRepository.Get(ctx, userID)
}

func (c *controller) Delete(ctx context.Context, userID int64) error {
	return c.usersRepository.Delete(ctx, userID)
}

func (c *controller) List(ctx context.Context) ([]User, error) {
	return c.usersRepository.List(ctx)
}

var errUnshureSuperuser = errors.New("failed to ensure superuser")

func (c *controller) enshureSuperuser() error {
	su, err := c.Get(context.Background(), SUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Join(errUnshureSuperuser, err)
		}
	}
	if su != nil {
		return nil
	}
	user := &User{
		Name:  "superuser",
		Roles: []int64{SUID},
	}
	if err := c.Create(context.Background(), user); err != nil {
		return errors.Join(errUnshureSuperuser, err)
	}
	c.suCreated = true
	return nil
}

func (c *controller) PrintToken() bool {
	return c.suCreated
}
