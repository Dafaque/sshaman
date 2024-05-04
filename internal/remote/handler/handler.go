package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dafaque/sshaman/internal/remote/controllers/users"
	remote "github.com/Dafaque/sshaman/pkg/remote/api"
)

type server struct {
	remote.UnimplementedRemoteCredentialsManagerServer
	usersController users.UsersController
}

func New(usersController users.UsersController) remote.RemoteCredentialsManagerServer {
	return &server{
		usersController: usersController,
	}
}

func (s *server) GetCredential(ctx context.Context, req *remote.GetCredentialRequest) (*remote.Credential, error) {
	// Implementation logic to retrieve a credential based on alias
	return nil, status.Errorf(codes.Unimplemented, "method GetCredential not implemented")
}

func (s *server) SetCredential(ctx context.Context, req *remote.SetCredentialRequest) (*remote.SetCredentialResponse, error) {
	// Implementation logic to add a new credential to the remote store
	return nil, status.Errorf(codes.Unimplemented, "method SetCredential not implemented")
}

func (s *server) ListCredentials(ctx context.Context, req *remote.ListCredentialsRequest) (*remote.ListCredentialsResponse, error) {
	// Implementation logic to list all credentials stored remotely
	return nil, status.Errorf(codes.Unimplemented, "method ListCredentials not implemented")
}

func (s *server) DeleteCredential(ctx context.Context, req *remote.DeleteCredentialRequest) (*remote.DeleteCredentialResponse, error) {
	// Implementation logic to remove a credential from the remote store
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCredential not implemented")
}

func (s *server) DropAllCredentials(ctx context.Context, req *remote.DropAllCredentialsRequest) (*remote.DropAllCredentialsResponse, error) {
	// Implementation logic to drop all credentials from the remote store
	return nil, status.Errorf(codes.Unimplemented, "method DropAllCredentials not implemented")
}

func (s *server) AddRole(ctx context.Context, req *remote.AddRoleRequest) (*remote.AddRoleResponse, error) {
	// Implementation logic to add a role
	return nil, status.Errorf(codes.Unimplemented, "method AddRole not implemented")
}

func (s *server) DeleteRole(ctx context.Context, req *remote.DeleteRoleRequest) (*remote.DeleteRoleResponse, error) {
	// Implementation logic to delete a role
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}

func (s *server) ListRoles(ctx context.Context, req *remote.ListRolesRequest) (*remote.ListRolesResponse, error) {
	// Implementation logic to list roles
	return nil, status.Errorf(codes.Unimplemented, "method ListRoles not implemented")
}
