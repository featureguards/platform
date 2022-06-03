package jwt_auth

import (
	"context"
	"strings"

	"platform/go/core/app_context"
	"platform/go/core/jwt"
	"platform/go/grpc/middleware/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Use Envoy's JWT token auth. This is ok for now for testing.

const (
	Key = "Authorization"
)

type Auth struct {
	allowedUnauthenticatedMethods map[string]struct{}
	j                             *jwt.JWT
}

type AuthOpts struct {
	AllowedUnauthenticatedMethods []string
	Jwt                           *jwt.JWT
}

func New(opts AuthOpts) (*Auth, error) {
	allowedUnauthenticatedMethods := make(map[string]struct{}, len(opts.AllowedUnauthenticatedMethods))
	for _, m := range opts.AllowedUnauthenticatedMethods {
		allowedUnauthenticatedMethods[m] = struct{}{}
	}
	return &Auth{j: opts.Jwt, allowedUnauthenticatedMethods: allowedUnauthenticatedMethods}, nil
}

func (a *Auth) authenticate(ctx context.Context, fullMethod string) (context.Context, error) {
	// If this method is expected not to require auth, let it through.
	if _, ok := a.allowedUnauthenticatedMethods[fullMethod]; ok {
		return ctx, nil
	}

	m := meta.ExtractIncoming(ctx)
	token := m.Get(Key)
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "No token")
	}

	t, err := a.j.ParseToken([]byte(token))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return app_context.WithJwtToken(ctx, t), nil
}

func (a *Auth) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := a.authenticate(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (a *Auth) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := a.authenticate(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
