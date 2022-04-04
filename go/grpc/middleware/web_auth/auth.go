package web_auth

import (
	"context"
	"net/http"

	"stackv2/go/core/app_context"
	"stackv2/go/core/ory"
	"stackv2/go/grpc/error_codes"
	"stackv2/go/grpc/middleware/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	kratos "github.com/ory/kratos-client-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cookieKey  = "cookie"
	sessionKey = "app.sid"
)

type Auth struct {
	client                        *kratos.APIClient
	allowedUnauthenticatedMethods map[string]struct{}
	allowedUnverifiedEmailMethods map[string]struct{}
}

type AuthOpts struct {
	KratosPublicURL               string
	AllowedUnverifiedEmailMethods []string
	AllowedUnauthenticatedMethods []string
}

func New(opts AuthOpts) (*Auth, error) {
	client, err := ory.New(ory.Opts{KratosPublicURL: opts.KratosPublicURL})
	if err != nil {
		return nil, err
	}

	allowedUnauthenticatedMethods := make(map[string]struct{}, len(opts.AllowedUnauthenticatedMethods))
	allowedUnverifiedEmailMethods := make(map[string]struct{}, len(opts.AllowedUnverifiedEmailMethods))
	for _, m := range opts.AllowedUnverifiedEmailMethods {
		allowedUnverifiedEmailMethods[m] = struct{}{}
	}
	for _, m := range opts.AllowedUnauthenticatedMethods {
		allowedUnauthenticatedMethods[m] = struct{}{}
	}
	return &Auth{client: client, allowedUnverifiedEmailMethods: allowedUnverifiedEmailMethods, allowedUnauthenticatedMethods: allowedUnauthenticatedMethods}, nil
}

func (a *Auth) authenticate(ctx context.Context, fullMethod string) (context.Context, error) {
	// If this method is expected not to require auth, let it through.
	if _, ok := a.allowedUnauthenticatedMethods[fullMethod]; ok {
		return ctx, nil
	}

	logger := log.WithContext(ctx)
	m := meta.ExtractIncoming(ctx)
	cookie := m.Get(cookieKey)
	if cookie == "" {
		return nil, status.Error(codes.Unauthenticated, "No cookie")
	}

	session, res, err := a.client.V0alpha2Api.ToSession(ctx).Cookie(cookie).Execute()
	if err != nil {
		logger.Warnf("Error Kratos session: %s", err)
		code := codes.Unknown
		if res != nil {
			code = error_codes.GrpcCode(int32(res.StatusCode))
		}
		return nil, status.Error(code, http.StatusText(res.StatusCode))
	}
	logger.Infof("%s %+v\n", fullMethod, session)
	if session.GetActive() {
		if ory.HasVerifiedAddress(session.Identity) {
			return app_context.WithSession(ctx, session), nil
		}
		// If this method is allowed to be called with unverified emails, let it through.
		if _, ok := a.allowedUnverifiedEmailMethods[fullMethod]; ok {
			return app_context.WithSession(ctx, session), nil
		}
	}
	return nil, status.Error(codes.Unauthenticated, "invalid session")
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
