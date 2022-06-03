package setup

import (
	"context"
	"net"
	"platform/go/grpc/auth"
	"platform/go/grpc/dashboard"
	"platform/go/grpc/toggles"
	pb_dashbard "platform/go/proto/dashboard"
	"platform/go/test/mocks/mock_app"
	"testing"

	pb_auth "github.com/featureguards/client-go/proto/auth"
	pb_toggles "github.com/featureguards/client-go/proto/toggles"
	"github.com/stretchr/testify/require"
)

type Apps struct {
	AuthListener    net.Listener
	AuthServer      *auth.AuthServer
	AuthClient      pb_auth.AuthClient
	DashListener    net.Listener
	DashboardServer *dashboard.DashboardServer
	DashboardClient pb_dashbard.DashboardClient
	TogglesListener net.Listener
	TogglesClient   pb_toggles.TogglesClient
	TogglesServer   *toggles.TogglesServer
	App             *mock_app.MockApp
}

func App(t *testing.T) *Apps {
	app, err := mock_app.New(t)
	require.Nil(t, err)
	t.Cleanup(app.Cleanup)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	authClient, authServer, authListen, err := authServer(ctx, t, app)
	require.Nil(t, err)

	dashClient, dashServer, dashListen, err := dashboardServer(ctx, t, app)
	require.Nil(t, err)

	togglesClient, togglesServer, togglesListen, err := togglesServer(ctx, t, app)
	require.Nil(t, err)

	return &Apps{
		App:             app,
		AuthServer:      authServer,
		AuthListener:    authListen,
		AuthClient:      authClient,
		DashboardServer: dashServer,
		DashListener:    dashListen,
		DashboardClient: dashClient,
		TogglesServer:   togglesServer,
		TogglesListener: togglesListen,
		TogglesClient:   togglesClient,
	}
}
