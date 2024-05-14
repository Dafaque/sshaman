package credentials

import (
	"context"
	"errors"
	"strings"

	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/internal/server/errs"
)

type Controller interface {
	Get(ctx context.Context, alias string) (credentials.Credentials, error)
	Set(ctx context.Context, cred credentials.Credentials, force bool) error
	List(ctx context.Context) ([]credentials.Credentials, error)
	Delete(ctx context.Context, alias string) error
	Drop(ctx context.Context) error
}

type controller struct {
	repo repository
}

func New() Controller {
	return &controller{}
}

func (c *controller) Get(ctx context.Context, alias string) (credentials.Credentials, error) {
	// Implementation needed
	return credentials.Credentials{}, errors.New("not implemented")
}

func (c *controller) Set(ctx context.Context, cred credentials.Credentials, force bool) error {
	perms, err := getPermissions(ctx)
	if err != nil {
		return err
	}
	repoCreds, err := prepareCredentials(&cred, perms)
	if err != nil {
		return err
	}

	var permitted bool = perms.SU()
	if !permitted {
		permitted = perms.Write(repoCreds.Space)
		if force {
			permitted = permitted && perms.Overwrite(repoCreds.Space)
		}
	}

	if !permitted {
		return errs.ErrNotPermitted
	}

	return nil
}

func (c *controller) List(ctx context.Context) ([]credentials.Credentials, error) {
	// Implementation needed
	return nil, errors.New("not implemented")
}

func (c *controller) Delete(ctx context.Context, alias string) error {
	// Implementation needed
	return errors.New("not implemented")
}

func (c *controller) Drop(ctx context.Context) error {
	// Implementation needed
	return errors.New("not implemented")
}

func parseSpace(alias string) (string, string, error) {
	parts := strings.Split(alias, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid alias")
	}
	return parts[0], parts[1], nil
}

func getPermissions(ctx context.Context) (permissions, error) {
	perms, ok := ctx.Value("permissions").(permissions)
	if !ok {
		return nil, errors.New("no permissions found")
	}
	return perms, nil
}

func prepareCredentials(cred *credentials.Credentials, perms permissions) (*Credentials, error) {
	space, alias, err := parseSpace(cred.Alias)
	if err != nil {
		return nil, err
	}
	newCred := &Credentials{
		Credentials: *cred,
		UserID:      perms.UID(),
	}
	newCred.Space = space
	newCred.Alias = alias
	return newCred, nil
}
