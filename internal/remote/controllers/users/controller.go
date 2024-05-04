package users

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"github.com/Dafaque/sshaman/internal/remote/auth"
)

type UsersController interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user User) error
	Get(ctx context.Context, userID int64) (*User, error)
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context) ([]User, error)
}

type usersController struct {
	jwtManager      *auth.JWTManager
	usersRepository usersRepository
	logger          *zap.Logger
}

func New(jwtManager *auth.JWTManager, usersRepository usersRepository, logger *zap.Logger) (UsersController, error) {
	controller := &usersController{
		jwtManager:      jwtManager,
		usersRepository: usersRepository,
		logger:          logger.Named("UsersControllers"),
	}
	err := controller.checkSuperuser()
	if err != nil {
		if errors.Is(err, errSuperuserNotFound) {
			createSuperuser := &User{
				Name:  "su",
				Roles: []string{"su"},
			}
			err = controller.Create(context.Background(), createSuperuser)
			if err != nil {
				return nil, err
			}
			str, err := controller.jwtManager.GenerateToken(createSuperuser.ID)
			if err != nil {
				return nil, err
			}
			controller.logger.Info("superuser created", zap.String("token", str))
		}
		return nil, err
	}
	return controller, nil
}

func (c *usersController) Create(ctx context.Context, user *User) error {
	return c.usersRepository.Create(ctx, user)
}

func (c *usersController) Update(ctx context.Context, user User) error {
	return c.usersRepository.Update(ctx, user)
}

func (c *usersController) Get(ctx context.Context, userID int64) (*User, error) {
	return c.usersRepository.Get(ctx, userID)
}

func (c *usersController) Delete(ctx context.Context, userID int64) error {
	return c.usersRepository.Delete(ctx, userID)
}

func (c *usersController) List(ctx context.Context) ([]User, error) {
	return c.usersRepository.List(ctx)
}

var errSuperuserNotFound = errors.New("superuser not found")

func (c *usersController) checkSuperuser() error {
	su, err := c.Get(context.Background(), 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errSuperuserNotFound
		}
		return err
	}
	if su == nil {
		return errSuperuserNotFound
	}
	return nil
}
