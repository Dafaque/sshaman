package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	remote "github.com/Dafaque/sshaman/pkg/remote/api"
)

func (s *server) CreateRole(ctx context.Context, req *remote.CreateRoleRequest) (*remote.CreateRoleResponse, error) {
	// Implementation logic to add a role
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}

func (s *server) UpdateRole(ctx context.Context, req *remote.UpdateRoleRequest) (*remote.UpdateRoleResponse, error) {
	// Implementation logic to update a role
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRole not implemented")
}

func (s *server) DeleteRole(ctx context.Context, req *remote.DeleteRoleRequest) (*remote.DeleteRoleResponse, error) {
	// Implementation logic to delete a role
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}

func (s *server) ListRoles(ctx context.Context, req *remote.ListRolesRequest) (*remote.ListRolesResponse, error) {
	// Implementation logic to list roles
	return nil, status.Errorf(codes.Unimplemented, "method ListRoles not implemented")
}
