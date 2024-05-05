package handler

import (
	"github.com/Dafaque/sshaman/internal/remote/controllers/roles"
	"github.com/Dafaque/sshaman/internal/remote/controllers/users"
	remote "github.com/Dafaque/sshaman/pkg/remote/api"
)

type server struct {
	remote.UnimplementedRemoteCredentialsManagerServer
	usersController users.Controller
	rolesController roles.Controller
}

func New(usersController users.Controller, rolesController roles.Controller) remote.RemoteCredentialsManagerServer {
	return &server{
		usersController: usersController,
		rolesController: rolesController,
	}
}
