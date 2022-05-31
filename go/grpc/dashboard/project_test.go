package dashboard_test

import (
	"context"

	"platform/go/grpc/middleware/meta"
	"platform/go/grpc/middleware/token_auth"
	"platform/go/test/setup"
	"testing"

	pb_dashboard "platform/go/proto/dashboard"

	"github.com/Pallinder/go-randomdata"
	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProjectTestSuite struct {
	suite.Suite
	apps    *setup.Apps
	token   string
	session *kratos.Session
}

func (suite *ProjectTestSuite) SetupSuite() {
	apps := setup.App(suite.T())
	suite.apps = apps
	ctx := context.Background()
	session, token, err := apps.CreateUserWithSession(ctx)
	require.Nil(suite.T(), err)
	suite.session = session
	suite.token = token

	// This lazily creates the user object
	res, err := suite.apps.DashboardClient.GetUser(suite.WithToken(ctx), &pb_dashboard.GetUserRequest{UserId: "me"})
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), res.OryId, session.Identity.Id)
}

func (suite *ProjectTestSuite) WithToken(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set(token_auth.Key, suite.token)
	return md.ToOutgoing(ctx)
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ProjectTestSuite) TestCreateProject() {
	ctx := context.Background()
	name := randomdata.Alphanumeric(10)
	proj, err := suite.apps.DashboardClient.CreateProject(suite.WithToken(ctx), &pb_dashboard.CreateProjectRequest{
		Name: name,
	})
	require.Nil(suite.T(), err, "%s", err)
	require.Equal(suite.T(), proj.Name, name)
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
