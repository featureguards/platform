package dashboard

import (
	"context"
	"fmt"
	"net"

	"platform/go/cmd"
	"platform/go/core/app"
	"platform/go/grpc/middleware"
	"platform/go/grpc/middleware/web_auth"
	"platform/go/grpc/middleware/web_log"
	"platform/go/grpc/server"
	pb_dashboard "platform/go/proto/dashboard"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Make sure DashboardServer implements everything
var _ pb_dashboard.DashboardServer = &DashboardServer{}

// server is used to implement dashboard.DashboardServer.
type DashboardServer struct {
	pb_dashboard.UnimplementedDashboardServer
	app app.App
}

func create(options ...server.ServerOptions) (*DashboardServer, error) {
	opts := &server.Options{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.App == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	return &DashboardServer{app: opts.App}, nil
}

func (s *DashboardServer) DB(ctx context.Context) *gorm.DB {
	return s.app.DB().WithContext(ctx)
}

type dashboardOptions struct {
	server.ListenerOptions
	auth middleware.Middleware
}

func WithPort(port int) DashboardOptions {
	return func(d *dashboardOptions) error {
		d.Port = port
		return nil
	}
}

func WithListener(l net.Listener) DashboardOptions {
	return func(d *dashboardOptions) error {
		d.Listener = l
		return nil
	}
}

func WithAuth(auth middleware.Middleware) DashboardOptions {
	return func(d *dashboardOptions) error {
		d.auth = auth
		return nil
	}
}

type DashboardOptions func(d *dashboardOptions) error

func Listen(ctx context.Context, a app.App, options ...DashboardOptions) (*DashboardServer, *grpc.Server, net.Listener, error) {
	lo := &dashboardOptions{}
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
	auth := lo.auth
	if auth == nil {
		var err error
		auth, err = web_auth.New(web_auth.AuthOpts{
			AllowedUnverifiedEmailMethods: []string{"/dashboard.Dashboard/GetUser"},
			KratosPublicURL:               a.Ory().PublicURL,
			KratosAdminURL:                a.Ory().AdminURL,
		})
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

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(recovery...), logger.StreamServerInterceptor(), auth.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recovery...), logger.UnaryServerInterceptor(), auth.UnaryServerInterceptor(),
		)),
	)

	dashboardServer, err := create(server.WithApp(a))
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	pb_dashboard.RegisterDashboardServer(srv, dashboardServer)
	return dashboardServer, srv, lis, nil
}

func (s *DashboardServer) HealthCheck(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
