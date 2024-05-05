package users

import "context"

type (
	usersRepository interface {
		Create(ctx context.Context, user *User) error
		Update(ctx context.Context, user User) error
		Get(ctx context.Context, userID int64) (*User, error)
		List(ctx context.Context) ([]User, error)
		Delete(ctx context.Context, userID int64) error
	}
	permissions interface {
		SU() bool
	}
)
