package handler

import (
	"github.com/Dafaque/sshaman/internal/server/controllers/roles"
	"github.com/Dafaque/sshaman/internal/server/controllers/users"
	api "github.com/Dafaque/sshaman/pkg/server/api"
)

type server struct {
	api.UnimplementedRemoteCredentialsManagerServer
	usersController users.Controller
	rolesController roles.Controller
}

func New(usersController users.Controller, rolesController roles.Controller) api.RemoteCredentialsManagerServer {
	return &server{
		usersController: usersController,
		rolesController: rolesController,
	}
}
