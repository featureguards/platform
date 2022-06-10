package toggles

import (
	"context"
	"fmt"
	"net"
	"time"

	"platform/go/cmd"
	"platform/go/core/app"
	"platform/go/core/app_context"
	cached_api_key "platform/go/core/cached/api_key"
	cached_feature_toggle "platform/go/core/cached/feature_toggle"
	"platform/go/core/ids"
	"platform/go/core/jwt"
	"platform/go/core/random"
	"platform/go/grpc/middleware/jwt_auth"
	"platform/go/grpc/middleware/web_log"
	"platform/go/grpc/server"

	pb_private "platform/go/proto/private"
	pb_project "platform/go/proto/project"

	pb_ft "github.com/featureguards/featureguards-go/proto/feature_toggle"
	pb_toggles "github.com/featureguards/featureguards-go/proto/toggles"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	jwtx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Make sure TogglesServer implements everything
var _ pb_toggles.TogglesServer = &TogglesServer{}

const (
	PollInterval = 5 * time.Second
)

// server is used to implement dashboard.DashboardServer.
type TogglesServer struct {
	pb_toggles.UnimplementedTogglesServer
	app app.App
	jwt *jwt.JWT
}

func (s *TogglesServer) apiKey(ctx context.Context, token jwtx.Token) (*pb_project.ApiKey, error) {
	apiKeyID := ids.ID(token.Subject())
	if err := apiKeyID.Validate(); err != nil {
		// Impossible. This is a bug.
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	return cached_api_key.Get(ctx, apiKeyID, s.app)
}

func (s *TogglesServer) Fetch(ctx context.Context, req *pb_toggles.FetchRequest) (*pb_toggles.FetchResponse, error) {
	token, _ := app_context.JwtTokenFromContext(ctx)
	if token == nil {
		return nil, status.Errorf(codes.Unauthenticated, "no token")
	}

	key, err := s.apiKey(ctx, token)
	if err != nil {
		return nil, err
	}

	envVersion, toggles, err := s.query(ctx, ids.ID(key.EnvironmentId), req.Version)
	if err != nil {
		return nil, err
	}
	return &pb_toggles.FetchResponse{FeatureToggles: toggles, Version: envVersion.Version}, nil
}

func (s *TogglesServer) query(ctx context.Context, envID ids.ID, startingVersion int64) (*pb_private.EnvironmentVersion, []*pb_ft.FeatureToggle, error) {
	envVersion, err := cached_feature_toggle.GetEnvironmentVersion(ctx, envID, s.app)
	if err != nil {
		return nil, nil, err
	}

	if envVersion.Version <= startingVersion {
		// Nothing more.
		return envVersion, nil, nil
	}

	envToggles, err := cached_feature_toggle.GetFeatureToggles(ctx, envID, s.app, startingVersion, envVersion.Version)
	if err != nil {
		return nil, nil, err
	}
	return envVersion, envToggles.FeatureToggles, nil
}

func (s *TogglesServer) Listen(req *pb_toggles.ListenRequest, stream pb_toggles.Toggles_ListenServer) error {
	ctx := stream.Context()
	token, _ := app_context.JwtTokenFromContext(ctx)
	if token == nil {
		return status.Errorf(codes.Unauthenticated, "no token")
	}

	key, err := s.apiKey(ctx, token)
	if err != nil {
		return err
	}
	envID := ids.ID(key.EnvironmentId)
	clientVersion := req.Version

	// Let's see if there is anything to send back initially
	envVersion, toggles, err := s.query(ctx, envID, clientVersion)
	if err != nil {
		return err
	}

	if envVersion.Version > clientVersion {
		if err := stream.Send(&pb_toggles.ListenPayload{
			FeatureToggles: toggles,
			Version:        envVersion.Version,
		}); err != nil {
			// If we can't send it to the client, it's likely they went away.
			return err
		}
		clientVersion = envVersion.Version
	}

	end := token.Expiration()
	exp := s.app.Clock().Timer(s.app.Clock().Until(end))
	defer exp.Stop()

	// TODO: Use pubsub too. Polling should be a last resort.
	ticker := s.app.Clock().Timer(random.Jitter(PollInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
		case <-exp.C:
			return status.Errorf(codes.Unauthenticated, "token expired")
		case <-ticker.C:
			ticker.Reset(random.Jitter(PollInterval))
			envVersion, toggles, err = s.query(ctx, envID, clientVersion)
			if err != nil {
				return err
			}

			if envVersion.Version > clientVersion {
				if err := stream.Send(&pb_toggles.ListenPayload{
					FeatureToggles: toggles,
					Version:        envVersion.Version,
				}); err != nil {
					// If we can't send it to the client, it's likely they went away.
					return err
				}
				clientVersion = envVersion.Version
			}
		}
	}
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
