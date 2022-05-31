package toggles

import (
	"context"
	"fmt"
	"log"
	"net"

	"platform/go/cmd"
	"platform/go/core/app"
	"platform/go/grpc/middleware/web_auth"
	"platform/go/grpc/middleware/web_log"
	"platform/go/grpc/server"

	pb_toggles "github.com/featureguards/client-go/proto/toggles"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"gorm.io/gorm"
)

// Make sure TogglesServer implements everything
var _ pb_toggles.TogglesServer = &TogglesServer{}

// server is used to implement dashboard.DashboardServer.
type TogglesServer struct {
	pb_toggles.UnimplementedTogglesServer
	app app.App
}

func New(options ...server.ServerOptions) (*TogglesServer, error) {
	opts := &server.Options{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.App == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	return &TogglesServer{app: opts.App}, nil
}

func (s *TogglesServer) DB(ctx context.Context) *gorm.DB {
	return s.app.DB().WithContext(ctx)
}

func Listen(ctx context.Context, a app.App, options ...server.ListenOptions) (*TogglesServer, *grpc.Server, net.Listener, error) {
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

	auth, err := web_auth.New(web_auth.AuthOpts{
		AllowedUnverifiedEmailMethods: []string{"/dashboard.Dashboard/GetUser"},
	})
	if err != nil {
		log.Fatal(err)
	}

	recovery := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(cmd.Recovery),
	}

	// Auth server doesn't need toggles middleware.
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(recovery...), logger.StreamServerInterceptor(), auth.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recovery...), logger.UnaryServerInterceptor(), auth.UnaryServerInterceptor(),
		)),
	)

	togglesServer, err := New(server.WithApp(a))
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	pb_toggles.RegisterTogglesServer(srv, togglesServer)

	return togglesServer, srv, lis, err
}
