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

func (suite *AuthTestSuite) SetupTest() {
	err := suite.stub.Create(context.Background())
	require.Nil(suite.T(), err, "%+v", err)
}

func (suite *AuthTestSuite) TestAuthenticate() {
	ctx := context.Background()
	res, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{
		Version: version,
	})
	require.Nil(suite.T(), err, "%+v", err)
	suite.T().Logf("%+v\n", res)
	require.NotEmpty(suite.T(), res.AccessToken)
	require.NotEmpty(suite.T(), res.RefreshToken)
}

func (suite *AuthTestSuite) TestRefresh() {
	ctx := context.Background()
	res, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{
		Version: version,
	})
	require.Nil(suite.T(), err, "%+v", err)
	suite.T().Logf("%+v\n", res)
	require.NotEmpty(suite.T(), res.AccessToken)
	require.NotEmpty(suite.T(), res.RefreshToken)

	_, err = suite.stub.App.AuthClient.Refresh(ctx, &pb_auth.RefreshRequest{RefreshToken: res.AccessToken})
	require.Contains(suite.T(), err.Error(), "invalid refresh token")

	refresh, err := suite.stub.App.AuthClient.Refresh(ctx, &pb_auth.RefreshRequest{RefreshToken: res.RefreshToken})
	require.Nil(suite.T(), err, "%+v", err)

	require.NotEmpty(suite.T(), refresh.AccessToken)
	require.NotEmpty(suite.T(), refresh.RefreshToken)
	require.NotEqual(suite.T(), refresh.AccessToken, res.AccessToken)
	require.NotEqual(suite.T(), refresh.RefreshToken, res.RefreshToken)

	// Refresh again using the old token to trigger a re-use
	_, err = suite.stub.App.AuthClient.Refresh(ctx, &pb_auth.RefreshRequest{RefreshToken: res.RefreshToken})
	require.Contains(suite.T(), err.Error(), "re-authenticate")

	// Using the refreshed token should also be blocked now
	_, err = suite.stub.App.AuthClient.Refresh(ctx, &pb_auth.RefreshRequest{RefreshToken: refresh.RefreshToken})
	require.NotNil(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "re-authenticate")

}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
