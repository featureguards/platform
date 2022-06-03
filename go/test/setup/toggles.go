package setup

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net"
	"platform/go/core/app"
	"platform/go/core/jwt"
	"platform/go/grpc/toggles"
	"testing"

	pb_toggles "github.com/featureguards/client-go/proto/toggles"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func togglesServer(ctx context.Context, t *testing.T, app app.App) (pb_toggles.TogglesClient, *toggles.TogglesServer, net.Listener, error) {
	togglesListen, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)
	j, err := jwt.New(jwt.WithKeyPair(privKey, &privKey.PublicKey))
	require.Nil(t, err)
	togglesServer, srv, _, err := toggles.Listen(ctx, toggles.WithListener(togglesListen), toggles.WithApp(app), toggles.WithJWT(j))
	require.Nil(t, err)

	go func() {
		if err := srv.Serve(togglesListen); err != nil {
			require.FailNowf(t, "%s", err.Error())
		}
	}()

	togglesConn, err := grpc.Dial(togglesListen.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	t.Cleanup(func() { togglesConn.Close() })
	togglesClient := pb_toggles.NewTogglesClient(togglesConn)
	return togglesClient, togglesServer, togglesListen, nil
}
