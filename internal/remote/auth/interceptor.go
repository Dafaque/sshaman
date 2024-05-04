package auth

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCAuthInterceptor struct {
	jwtManager *JWTManager
}

func NewGRPCAuthInterceptor(jwtManager *JWTManager) *GRPCAuthInterceptor {
	return &GRPCAuthInterceptor{jwtManager: jwtManager}
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
		ctx = context.WithValue(ctx, "userID", uid)
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
		uid, err := interceptor.authorize(ss.Context())
		if err != nil {
			return err
		}
		ss.SetHeader(metadata.Pairs("userID", strconv.FormatInt(uid, 10)))
		return handler(srv, ss)
	}
}

func (interceptor *GRPCAuthInterceptor) authorize(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return -1, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	params, err := interceptor.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return -1, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	if !params.valid {
		return -1, status.Errorf(codes.Unauthenticated, "access token is invalid")
	}
	return params.userID, nil
}

func (interceptor *GRPCAuthInterceptor) Shutdown() error {
	return nil
}
