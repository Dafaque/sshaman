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
	usersController users.Controller
}

func New(usersController users.Controller) remote.RemoteCredentialsManagerServer {
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
