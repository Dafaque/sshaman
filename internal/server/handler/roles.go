package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dafaque/sshaman/internal/server/controllers/roles"
	"github.com/Dafaque/sshaman/internal/server/errs"
	api "github.com/Dafaque/sshaman/pkg/server/api"
)

func (s *server) CreateRole(ctx context.Context, req *api.CreateRoleRequest) (*api.CreateRoleResponse, error) {
	role := &roles.Role{
		Name:        req.Role.Name,
		Description: req.Role.Description,
		Read:        req.Role.Read,
		Write:       req.Role.Write,
		Delete:      req.Role.Delete,
		Overwrite:   req.Role.Overwrite,
		SU:          req.Role.Su,
		Spaces:      req.Role.Spaces,
	}
	err := s.rolesController.Create(ctx, role)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var success bool
	if role.ID > 0 {
		success = true
	}
	return &api.CreateRoleResponse{
		Success: success,
	}, nil
}

func (s *server) UpdateRole(ctx context.Context, req *api.UpdateRoleRequest) (*api.UpdateRoleResponse, error) {
	// Implementation logic to update a role
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRole not implemented")
}

func (s *server) DeleteRole(ctx context.Context, req *api.DeleteRoleRequest) (*api.DeleteRoleResponse, error) {
	// Implementation logic to delete a role
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}

func (s *server) ListRoles(ctx context.Context, req *api.ListRolesRequest) (*api.ListRolesResponse, error) {
	roles, err := s.rolesController.List(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var resp api.ListRolesResponse

	for _, role := range roles {
		resp.Roles = append(resp.Roles, &api.Role{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Read:        role.Read,
			Write:       role.Write,
			Delete:      role.Delete,
			Overwrite:   role.Overwrite,
			Su:          role.SU,
			Spaces:      role.Spaces,
		})
	}
	return &resp, nil
}
