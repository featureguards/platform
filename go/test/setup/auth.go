package setup

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net"
	"platform/go/core/app"
	"platform/go/core/jwt"
	"platform/go/grpc/auth"
	"testing"

	pb_auth "github.com/featureguards/client-go/proto/auth"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func authServer(ctx context.Context, t *testing.T, app app.App) (pb_auth.AuthClient, *auth.AuthServer, net.Listener, error) {
	authListen, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)
	j, err := jwt.New(jwt.WithKeyPair(privKey, &privKey.PublicKey))
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
