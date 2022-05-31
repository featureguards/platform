package auth

import (
	"context"
	"fmt"
	"net"
	"time"

	"platform/go/cmd"
	"platform/go/core/app"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/api_keys"
	"platform/go/grpc/middleware/meta"
	"platform/go/grpc/middleware/web_log"
	"platform/go/grpc/server"

	pb_project "platform/go/proto/project"

	pb_auth "github.com/featureguards/client-go/proto/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	ApiKeyMD = "x-api-key"
)

// Make sure AuthServer implements everything
var _ pb_auth.AuthServer = &AuthServer{}

// server is used to implement dashboard.DashboardServer.
type AuthServer struct {
	pb_auth.UnimplementedAuthServer
	app app.App
}

type AuthOpts struct {
	App app.App
}

func New(options ...server.ServerOptions) (*AuthServer, error) {
	opts := &server.Options{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.App == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	return &AuthServer{app: opts.App}, nil
}

func (s *AuthServer) Authenticate(ctx context.Context, req *pb_auth.AuthenticateRequest) (*pb_auth.AuthenticateResponse, error) {
	m := meta.ExtractIncoming(ctx)
	apiKeyHeader := m.Get(ApiKeyMD)
	if apiKeyHeader == "" {
		return nil, status.Errorf(codes.NotFound, "api-key header not specified")
	}

	key := ids.ID(apiKeyHeader)
	if err := key.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid api-key")
	}

	if ot, _, _ := ids.Parse(key); ot != ids.ApiKey {
		return nil, status.Errorf(codes.InvalidArgument, "invalid api-key")
	}

	pb, err := s.app.KV().GetProto(ctx, kv.ApiKey, string(key))
	if err != nil {
		// Fetch it from the database
		model, err := api_keys.Get(ctx, key, s.app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot find api key")
			}
			return nil, status.Error(codes.Internal, "invalid api key")
		}
		pb, err = api_keys.Pb(model)
		if err != nil {
			return nil, status.Error(codes.Internal, "invalid api key")
		}
		// Populate the cache
		if err := s.app.KV().SetProto(ctx, kv.ApiKey, string(key), pb); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	// pb must be set
	apiKey := pb.(*pb_project.ApiKey)
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.AsTime().Before(time.Now()) {
		return nil, status.Error(codes.Unauthenticated, "expired API key")
	}
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}

func (s *AuthServer) Refresh(context.Context, *pb_auth.RefreshRequest) (*pb_auth.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

func (s *AuthServer) DB(ctx context.Context) *gorm.DB {
	return s.app.DB().WithContext(ctx)
}

func Listen(ctx context.Context, a app.App, options ...server.ListenOptions) (*AuthServer, *grpc.Server, net.Listener, error) {
	lo := &server.ListenerOptions{}
	for _, opt := range options {
		opt(lo)
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

	authServer, err := New(server.WithApp(a))
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	pb_auth.RegisterAuthServer(srv, authServer)
	return authServer, srv, lis, nil
}
