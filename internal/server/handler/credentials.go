package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/Dafaque/sshaman/pkg/server/api"
)

func (s *server) GetCredential(ctx context.Context, req *api.GetCredentialRequest) (*api.Credential, error) {
	// Implementation logic to retrieve a credential based on alias
	return nil, status.Errorf(codes.Unimplemented, "method GetCredential not implemented")
}

func (s *server) SetCredential(ctx context.Context, req *api.SetCredentialRequest) (*api.SetCredentialResponse, error) {
	// Implementation logic to add a new credential to the remote store
	return nil, status.Errorf(codes.Unimplemented, "method SetCredential not implemented")
}

func (s *server) ListCredentials(ctx context.Context, req *api.ListCredentialsRequest) (*api.ListCredentialsResponse, error) {
	// Implementation logic to list all credentials stored remotely
	return nil, status.Errorf(codes.Unimplemented, "method ListCredentials not implemented")
}

func (s *server) DeleteCredential(ctx context.Context, req *api.DeleteCredentialRequest) (*api.DeleteCredentialResponse, error) {
	// Implementation logic to remove a credential from the remote store
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCredential not implemented")
}

func (s *server) DropAllCredentials(ctx context.Context, req *api.DropAllCredentialsRequest) (*api.DropAllCredentialsResponse, error) {
	// Implementation logic to drop all credentials from the remote store
	return nil, status.Errorf(codes.Unimplemented, "method DropAllCredentials not implemented")
}
