package users

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"github.com/Dafaque/sshaman/internal/remote/auth"
)

type Controller interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user User) error
	Get(ctx context.Context, userID int64) (*User, error)
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context) ([]User, error)
}

type controller struct {
	usersRepository usersRepository
}

func New(jwtManager *auth.JWTManager, usersRepository usersRepository, logger *zap.Logger) (Controller, error) {
	controller := &controller{
		usersRepository: usersRepository,
	}
	err := controller.enshureSuperuser(jwtManager, logger)
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

func (rc *controller) enshureSuperuser(jwtManager *auth.JWTManager, logger *zap.Logger) error {
	su, err := rc.Get(context.Background(), 0)
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
		Roles: []string{"su"},
	}
	if err := rc.Create(context.Background(), user); err != nil {
		return errors.Join(errUnshureSuperuser, err)
	}
	str, err := jwtManager.GenerateToken(user.ID)
	if err != nil {
		return errors.Join(errUnshureSuperuser, err)
	}
	logger.Named("enshureSuperuser").Info("superuser created", zap.String("token", str))
	return nil
}
