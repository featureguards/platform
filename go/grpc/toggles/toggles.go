package toggles

import (
	"context"
	"fmt"
	"log"
	"net"

	"platform/go/cmd"
	"platform/go/core/app"
	"platform/go/core/jwt"
	"platform/go/grpc/middleware/jwt_auth"
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
	jwt *jwt.JWT
}

func WithPort(port int) TogglesOptions {
	return func(o *togglesOptions) error {
		o.Port = port
		return nil
	}
}

func WithListener(l net.Listener) TogglesOptions {
	return func(o *togglesOptions) error {
		o.Listener = l
		return nil
	}
}

func WithApp(a app.App) TogglesOptions {
	return func(o *togglesOptions) error {
		o.app = a
		return nil
	}
}

func WithJWT(j *jwt.JWT) TogglesOptions {
	return func(o *togglesOptions) error {
		o.jwt = j
		return nil
	}
}

type togglesOptions struct {
	server.ListenerOptions
	jwt *jwt.JWT
	app app.App
}

type TogglesOptions func(d *togglesOptions) error

func create(options ...TogglesOptions) (*TogglesServer, error) {
	opts := &togglesOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.app == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	return &TogglesServer{app: opts.app, jwt: opts.jwt}, nil
}

func (s *TogglesServer) DB(ctx context.Context) *gorm.DB {
	return s.app.DB().WithContext(ctx)
}

func Listen(ctx context.Context, options ...TogglesOptions) (*TogglesServer, *grpc.Server, net.Listener, error) {
	o := &togglesOptions{}
	for _, opt := range options {
		opt(o)
	}
	lis := o.Listener
	if lis == nil {
		var err error
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.Port))
		if err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}
	}
	logger, err := web_log.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	auth, err := jwt_auth.New(jwt_auth.AuthOpts{
		Jwt: o.jwt,
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

	togglesServer, err := create(WithApp(o.app), WithJWT(o.jwt))
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
