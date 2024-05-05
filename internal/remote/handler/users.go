package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dafaque/sshaman/internal/remote/controllers/users"
	remote "github.com/Dafaque/sshaman/pkg/remote/api"
)

func (s *server) CreateUser(ctx context.Context, req *remote.CreateUserRequest) (*remote.CreateUserResponse, error) {
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
	if user.ID >= 0 {
		success = false
	}
	resp := &remote.CreateUserResponse{
		Success: success,
	}

	return resp, nil
}

func (s *server) UpdateUser(ctx context.Context, req *remote.UpdateUserRequest) (*remote.UpdateUserResponse, error) {
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
	resp := &remote.UpdateUserResponse{
		Success: true,
	}
	return resp, nil
}

func (s *server) DeleteUser(ctx context.Context, req *remote.DeleteUserRequest) (*remote.DeleteUserResponse, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete superuser")
	}
	err := s.usersController.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &remote.DeleteUserResponse{
		Success: true,
	}
	return resp, nil
}

func (s *server) ListUsers(ctx context.Context, req *remote.ListUsersRequest) (*remote.ListUsersResponse, error) {
	users, err := s.usersController.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &remote.ListUsersResponse{
		Users: make([]*remote.User, len(users)),
	}

	for i, user := range users {
		resp.Users[i] = &remote.User{
			Id:    user.ID,
			Name:  user.Name,
			Roles: user.Roles,
		}
	}
	return resp, nil
}
