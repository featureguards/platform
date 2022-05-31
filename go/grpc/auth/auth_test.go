package auth_test

import (
	"context"
	"platform/go/test/stubs"
	"testing"

	pb_auth "github.com/featureguards/client-go/proto/auth"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	version = "0.1.0"
)

type AuthTestSuite struct {
	suite.Suite
	stub *stubs.Stubs
}

func (suite *AuthTestSuite) SetupSuite() {
	ctx := context.Background()
	stub := stubs.New(ctx, suite.T())
	suite.stub = stub
}

func (suite *AuthTestSuite) TestAuthenticate() {
	ctx := context.Background()
	res, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{
		Version: version,
	})
	require.Nil(suite.T(), err, "%+v", err)
	require.NotEmpty(suite.T(), res.AccessToken)
	require.NotEmpty(suite.T(), res.RefreshToken)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
