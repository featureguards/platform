package setup

import (
	"context"
	"net"
	"platform/go/core/app"
	"platform/go/grpc/dashboard"
	"platform/go/grpc/middleware/token_auth"
	pb_dashbard "platform/go/proto/dashboard"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dashboardServer(ctx context.Context, t *testing.T, app app.App) (pb_dashbard.DashboardClient, *dashboard.DashboardServer, net.Listener, error) {
	dashListen, err := net.Listen("tcp", ":0")
	require.Nil(t, err)

	// We use session-token auth with Ory because we don't have a browser.
	tokenAuth, err := token_auth.New(token_auth.AuthOpts{
		KratosPublicURL:           app.Ory().PublicURL,
		KratosAdminURL:            app.Ory().AdminURL,
		TestAllowUnverifiedEmails: true,
	})
	require.Nil(t, err)

	dashServer, dashSrv, _, err := dashboard.Listen(ctx, app, dashboard.WithListener(dashListen), dashboard.WithAuth(tokenAuth))
	require.Nil(t, err)
	go func() {
		if err := dashSrv.Serve(dashListen); err != nil {
			require.FailNowf(t, "%s", err.Error())
		}
	}()

	dashConn, err := grpc.Dial(dashListen.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	t.Cleanup(func() { dashConn.Close() })
	dashClient := pb_dashbard.NewDashboardClient(dashConn)
	return dashClient, dashServer, dashListen, nil
}
