package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dafaque/sshaman/internal/server/controllers/users"
	api "github.com/Dafaque/sshaman/pkg/server/api"
)

func (s *server) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	user := &users.User{
		Name:  req.User.Name,
		Roles: req.User.Roles,
	}
	// @todo enshure roles exist
	err := s.usersController.Create(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	success := true
	if user.ID <= 0 {
		success = false
	}
	resp := &api.CreateUserResponse{
		Success: success,
	}

	return resp, nil
}

func (s *server) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.UpdateUserResponse, error) {
	user := users.User{
		ID:    req.User.Id,
		Name:  req.User.Name,
		Roles: req.User.Roles,
	}
	// @todo enshure roles exist
	err := s.usersController.Update(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &api.UpdateUserResponse{
		Success: true,
	}
	return resp, nil
}

func (s *server) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete superuser")
	}
	err := s.usersController.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &api.DeleteUserResponse{
		Success: true,
	}
	return resp, nil
}

func (s *server) ListUsers(ctx context.Context, req *api.ListUsersRequest) (*api.ListUsersResponse, error) {
	users, err := s.usersController.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &api.ListUsersResponse{
		Users: make([]*api.User, len(users)),
	}

	for i, user := range users {
		resp.Users[i] = &api.User{
			Id:    user.ID,
			Name:  user.Name,
			Roles: user.Roles,
		}
	}
	return resp, nil
}
