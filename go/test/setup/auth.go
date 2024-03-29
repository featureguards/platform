package setup

import (
	"context"
	"net"
	"platform/go/core/app"
	"platform/go/core/jwt"
	"platform/go/grpc/auth"
	"testing"

	pb_auth "github.com/featureguards/featureguards-go/v2/proto/auth"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func authServer(ctx context.Context, t *testing.T, app app.App, j *jwt.JWT) (pb_auth.AuthClient, *auth.AuthServer, net.Listener, error) {
	authListen, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	authServer, srv, _, err := auth.Listen(ctx, auth.WithListener(authListen), auth.WithApp(app), auth.WithJWT(j))
	require.Nil(t, err)

	go func() {
		if err := srv.Serve(authListen); err != nil {
			require.FailNowf(t, "%s", err.Error())
		}
	}()

	authConn, err := grpc.Dial(authListen.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	t.Cleanup(func() { authConn.Close() })
	authClient := pb_auth.NewAuthClient(authConn)
	return authClient, authServer, authListen, nil
}
