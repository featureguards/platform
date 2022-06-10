package auth

import (
	"context"
	"fmt"
	"net"
	"strings"

	"platform/go/cmd"
	"platform/go/core/app"
	cached_api_key "platform/go/core/cached/api_key"
	"platform/go/core/ids"
	"platform/go/core/jwt"
	"platform/go/core/kv"
	"platform/go/grpc/middleware/meta"
	"platform/go/grpc/middleware/web_log"
	"platform/go/grpc/server"

	pb_auth "github.com/featureguards/featureguards-go/v2/proto/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ApiKeyMD = "x-api-key"
)

// Make sure AuthServer implements everything
var _ pb_auth.AuthServer = &AuthServer{}

// server is used to implement auth.AuthServer.
type AuthServer struct {
	pb_auth.UnimplementedAuthServer
	app app.App
	jwt *jwt.JWT
}

func WithPort(port int) AuthOptions {
	return func(o *authOptions) error {
		o.Port = port
		return nil
	}
}

func WithListener(l net.Listener) AuthOptions {
	return func(o *authOptions) error {
		o.Listener = l
		return nil
	}
}

func WithApp(a app.App) AuthOptions {
	return func(o *authOptions) error {
		o.app = a
		return nil
	}
}

func WithJWT(j *jwt.JWT) AuthOptions {
	return func(o *authOptions) error {
		o.jwt = j
		return nil
	}
}

type authOptions struct {
	server.ListenerOptions
	jwt *jwt.JWT
	app app.App
}

type AuthOptions func(d *authOptions) error

func create(options ...AuthOptions) (*AuthServer, error) {
	opts := &authOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.app == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	if opts.jwt == nil {
		return nil, fmt.Errorf("jwt must be specified")
	}
	return &AuthServer{jwt: opts.jwt, app: opts.app}, nil
}

func (s *AuthServer) Authenticate(ctx context.Context, req *pb_auth.AuthenticateRequest) (*pb_auth.AuthenticateResponse, error) {
	m := meta.ExtractIncoming(ctx)
	apiKeyHeader := m.Get(ApiKeyMD)
	if apiKeyHeader == "" {
		return nil, status.Errorf(codes.NotFound, "api-key header not specified")
	}

	trimmed := strings.TrimSpace(apiKeyHeader)
	splitted := strings.Split(trimmed, ":")
	if len(splitted) != 2 {
		return nil, status.Error(codes.InvalidArgument, "invalid api-key format")
	}

	apiKeyID := ids.ID(splitted[0])

	if err := apiKeyID.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid api-key")
	}

	if ot, _, _ := ids.Parse(apiKeyID); ot != ids.ApiKey {
		return nil, status.Errorf(codes.InvalidArgument, "invalid api-key")
	}
	apiKey, err := cached_api_key.Get(ctx, apiKeyID, s.app)
	if err != nil {
		log.Errorf("%s\n", err)
		return nil, err
	}

	// The provided key and the key on the ApiKey must match
	if strings.Compare(apiKey.Key, trimmed) != 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid API key")
	}
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.AsTime().Before(s.app.Clock().Now()) {
		return nil, status.Error(codes.Unauthenticated, "expired API key")
	}

	// We're good now. Generate the tokens
	access, refresh, err := s.generateTokens(ctx, apiKeyID)
	if err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Errorf(codes.Internal, "could not create access token")
	}
	return &pb_auth.AuthenticateResponse{AccessToken: string(access), RefreshToken: string(refresh)}, nil
}

func (s *AuthServer) generateTokens(ctx context.Context, apiKeyID ids.ID, options ...jwt.TokenOptions) (access []byte, refresh []byte, err error) {
	access, err = s.jwt.SignedToken(apiKeyID, jwt.AccessToken)
	if err != nil {
		return
	}
	// We have the refresh token use the same JwtID as the access token because we bundle them all
	// together.
	refresh, err = s.jwt.SignedToken(apiKeyID, jwt.RefreshToken, options...)
	if err != nil {
		return
	}
	return
}

func (s *AuthServer) Refresh(ctx context.Context, req *pb_auth.RefreshRequest) (*pb_auth.RefreshResponse, error) {
	// ParseToken internally validates the notbefore and after timestamps, signature...etc.
	t, err := s.jwt.ParseToken([]byte(req.RefreshToken))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid refresh token")
	}

	// Only allow using refresh tokens to refresh.
	if len(t.Audience()) != 1 || !strings.EqualFold(t.Audience()[0], string(jwt.RefreshToken)) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid refresh token")
	}

	ttl := s.app.Clock().Until(t.Expiration())
	// We implement token rotation to detect token leakage.
	// See https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation
	set, err := s.app.KV().SetNX(ctx, kv.RefreshToken, t.JwtID(), []byte(t.Subject()), kv.WithExpiration(ttl))
	if err != nil {
		log.Errorf("%s\n", errors.WithStack(err))
		return nil, status.Errorf(codes.Internal, "could not refresh token")
	}
	if !set {
		// Means there is token re-use. Invalidate the refresh token family and require re-auth.
		if _, err := s.app.KV().SetNX(ctx, kv.RefreshTokenFamily, t.PrivateClaims()[jwt.TokenFamilyClaim].(string), []byte(t.Subject()), kv.WithExpiration(ttl)); err != nil {
			log.Errorf("%s\n", errors.WithStack(err))
			return nil, status.Errorf(codes.Internal, "could not refresh toekn")
		}
		return nil, status.Errorf(codes.PermissionDenied, "re-authenticate")
	}
	// Check to make sure the refresh token family isn't invalid.
	if _, err := s.app.KV().Get(ctx, kv.RefreshTokenFamily, t.PrivateClaims()[jwt.TokenFamilyClaim].(string)); err == nil || err != kv.ErrNotFound {
		return nil, status.Errorf(codes.PermissionDenied, "re-authenticate")
	}

	apiKeyID := ids.ID(t.Subject())
	if err := apiKeyID.Validate(); err != nil {
		// Impossible. This is a bug.
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	access, refresh, err := s.generateTokens(ctx, apiKeyID, jwt.WithFamily(t.PrivateClaims()[jwt.TokenFamilyClaim].(string)))
	if err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Errorf(codes.Internal, "could not create access token")
	}

	return &pb_auth.RefreshResponse{AccessToken: string(access), RefreshToken: string(refresh)}, nil
}

func Listen(ctx context.Context, options ...AuthOptions) (*AuthServer, *grpc.Server, net.Listener, error) {
	lo := &authOptions{}
	for _, opt := range options {
		opt(lo)
	}
	authServer, err := create(WithApp(lo.app), WithJWT(lo.jwt))
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	lis := lo.Listener
	if lis == nil {
		var err error
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", lo.Port))
		if err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}
	}
	logger, err := web_log.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	recovery := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(cmd.Recovery),
	}

	// Auth server doesn't need auth middleware.
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(recovery...), logger.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recovery...), logger.UnaryServerInterceptor(),
		)),
	)

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	pb_auth.RegisterAuthServer(srv, authServer)
	return authServer, srv, lis, nil
}
