package auth

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Dafaque/sshaman/internal/server/controllers/roles"
	"github.com/Dafaque/sshaman/internal/server/controllers/users"
)

type GRPCAuthInterceptor struct {
	jwtManager *JWTManager
	roles      roles.Controller
	users      users.Controller
}

func NewGRPCAuthInterceptor(
	jwtManager *JWTManager,
	roles roles.Controller,
	users users.Controller,
) *GRPCAuthInterceptor {
	return &GRPCAuthInterceptor{jwtManager: jwtManager, roles: roles, users: users}
}

func (interceptor *GRPCAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		uid, err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "permissions", uid)
		return handler(ctx, req)
	}
}

func (interceptor *GRPCAuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		perms, err := interceptor.authorize(ss.Context())
		if err != nil {
			return err
		}
		if perms.su {
			ss.SetHeader(metadata.Pairs("x-permissions-su", strconv.FormatBool(perms.su)))
		} else {
			ss.SetHeader(metadata.Pairs("x-permissions-read", strings.Join(perms.read, ",")))
			ss.SetHeader(metadata.Pairs("x-permissions-write", strings.Join(perms.write, ",")))
			ss.SetHeader(metadata.Pairs("x-permissions-delete", strings.Join(perms.delete, ",")))
			ss.SetHeader(metadata.Pairs("x-permissions-overwrite", strings.Join(perms.overwrite, ",")))
		}

		return handler(srv, ss)
	}
}

const tokenPrefix = "Bearer "

func (interceptor *GRPCAuthInterceptor) authorize(ctx context.Context) (*permissions, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]

	if len(accessToken) < len(tokenPrefix) {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken = strings.TrimPrefix(accessToken, tokenPrefix)

	params, err := interceptor.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	if !params.valid {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid")
	}

	user, err := interceptor.users.Get(ctx, params.userID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}
	roles, err := interceptor.roles.Get(ctx, user.Roles...)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "roles not found")
	}
	var perms permissions
	perms.uid = user.ID

	for _, role := range roles {
		if role.SU {
			perms.su = true
			break
		}
		if role.Read {
			perms.read = append(perms.read, role.Spaces...)
		}
		if role.Write {
			perms.write = append(perms.write, role.Spaces...)
		}
		if role.Delete {
			perms.delete = append(perms.delete, role.Spaces...)
		}
		if role.Overwrite {
			perms.overwrite = append(perms.overwrite, role.Spaces...)
		}
	}

	return &perms, nil
}
