package toggles_test

import (
	"context"
	"fmt"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/grpc/toggles"
	"platform/go/test/stubs"
	"sync/atomic"
	"testing"
	"time"

	pb_private "platform/go/proto/private"
	pb_project "platform/go/proto/project"

	"github.com/Pallinder/go-randomdata"
	"github.com/benbjohnson/clock"
	pb_auth "github.com/featureguards/featureguards-go/v2/proto/auth"
	pb_toggles "github.com/featureguards/featureguards-go/v2/proto/toggles"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	version = "0.1.0"
)

type TogglesTestSuite struct {
	suite.Suite
	stub *stubs.Stubs
}

func (suite *TogglesTestSuite) SetupSuite() {
	ctx := context.Background()
	stub := stubs.New(ctx, suite.T())
	suite.stub = stub
}

func (suite *TogglesTestSuite) SetupTest() {
	ctx := context.Background()
	err := suite.stub.Create(ctx)
	require.Nil(suite.T(), err, "%+v", err)
}

func (suite *TogglesTestSuite) TestFetch() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	parsed, err := suite.stub.App.Jwt.ParseToken([]byte(token.AccessToken))
	require.Nil(suite.T(), err, "%+v", err)

	res, err := suite.stub.App.TogglesClient.Fetch(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.FetchRequest{})
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), res.Version, int64(0))
	require.NotEmpty(suite.T(), res.FeatureToggles)
	require.Greater(suite.T(), res.SettingsVersion, int64(0))
	require.NotEmpty(suite.T(), res.DynamicSettings)

	// Ensure caches are populated
	cached, err := suite.stub.App.App.KV().GetProto(ctx, kv.ApiKey, parsed.Subject())
	require.Nil(suite.T(), err, "%+v", err)
	apiKey := cached.(*pb_project.ApiKey)
	require.NotNil(suite.T(), apiKey)

	// Feature toggles
	cached, err = suite.stub.App.App.KV().GetProto(ctx, kv.EnvironmentVersion, apiKey.EnvironmentId)
	require.Nil(suite.T(), err, "%+v", err)
	envVersion := cached.(*pb_private.EnvironmentVersion)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), envVersion.Id, apiKey.EnvironmentId)
	require.Equal(suite.T(), envVersion.Version, res.Version)

	cached, err = suite.stub.App.App.KV().GetProto(ctx, kv.EnvironmentToggles, apiKey.EnvironmentId, kv.WithSuffix(fmt.Sprintf("%d-%d", 0, res.Version)))
	require.Nil(suite.T(), err, "%+v", err)
	envToggles := cached.(*pb_private.EnvironmentFeatureToggles)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), envToggles.FeatureToggles, res.FeatureToggles)
	require.Equal(suite.T(), envToggles.StartingVersion, int64(0))
	require.Equal(suite.T(), envToggles.EndingVersion, envVersion.Version)

	// Dynamic settings
	cached, err = suite.stub.App.App.KV().GetProto(ctx, kv.EnvironmentSettingsVersion, apiKey.EnvironmentId)
	require.Nil(suite.T(), err, "%+v", err)
	envSettingsVersion := cached.(*pb_private.EnvironmentVersion)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), envSettingsVersion.Id, apiKey.EnvironmentId)
	require.Equal(suite.T(), envSettingsVersion.Version, res.Version)

	cached, err = suite.stub.App.App.KV().GetProto(ctx, kv.EnvironmentSettings, apiKey.EnvironmentId, kv.WithSuffix(fmt.Sprintf("%d-%d", 0, res.SettingsVersion)))
	require.Nil(suite.T(), err, "%+v", err)
	envSettings := cached.(*pb_private.EnvironmentDynamicSettings)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), envSettings.DynamicSettings, res.DynamicSettings)
	require.Equal(suite.T(), envSettings.StartingVersion, int64(0))
	require.Equal(suite.T(), envSettings.EndingVersion, envSettingsVersion.Version)

	// Fetch again using the same version we got earlier.
	res2, err := suite.stub.App.TogglesClient.Fetch(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.FetchRequest{Version: res.Version, SettingsVersion: res.SettingsVersion})
	require.Nil(suite.T(), err, "%+v", err)
	require.Equal(suite.T(), res.Version, res2.Version)
	require.Empty(suite.T(), res2.FeatureToggles)
	require.Equal(suite.T(), res.SettingsVersion, res2.SettingsVersion)
	require.Empty(suite.T(), res2.DynamicSettings)

	// Let's exercise the cache by querying the same data
	res, err = suite.stub.App.TogglesClient.Fetch(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.FetchRequest{})
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), res.Version, int64(0))
	require.NotEmpty(suite.T(), res.FeatureToggles)
	require.Greater(suite.T(), res.SettingsVersion, int64(0))
	require.NotEmpty(suite.T(), res.DynamicSettings)
}

func (suite *TogglesTestSuite) TestBasicListen() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1, SettingsVersion: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload.Version, int64(0))
	require.Greater(suite.T(), payload.SettingsVersion, int64(0))
}

func (suite *TogglesTestSuite) TestListenBlocks() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1, SettingsVersion: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload.Version, int64(0))
	require.Greater(suite.T(), payload.SettingsVersion, int64(0))

	// Second one must block
	gotPayload := int32(0)
	go func() {
		client.Recv()
		atomic.StoreInt32(&gotPayload, 1)
	}()
	cl := suite.stub.App.App.Clock().(*clock.Mock)
	cl.Add(3 * time.Second)
	<-time.After(100 * time.Millisecond)
	require.Equal(suite.T(), gotPayload, int32(0))

}

func (suite *TogglesTestSuite) TestServerReauthaneticates() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1, SettingsVersion: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload.Version, int64(0))
	require.Greater(suite.T(), payload.SettingsVersion, int64(0))

	parsed, err := suite.stub.App.Jwt.ParseToken([]byte(token.AccessToken))
	require.Nil(suite.T(), err, "%+v", err)
	clock := suite.stub.App.App.Clock().(*clock.Mock)
	clock.Add(clock.Until(parsed.Expiration()))
	_, err = client.Recv()
	require.Contains(suite.T(), err.Error(), "code = Unauthenticated")
}

func (suite *TogglesTestSuite) TestClientDeadlineExceeded() {
	exp := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload.Version, int64(0))

	clock := suite.stub.App.App.Clock().(*clock.Mock)
	clock.Add(exp)
	_, err = client.Recv()
	require.Contains(suite.T(), err.Error(), "code = DeadlineExceeded")
}

func (suite *TogglesTestSuite) TestListenWithCreateFeatureToggle() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload1, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload1.Version, int64(0))

	clock := suite.stub.App.App.Clock().(*clock.Mock)
	// Let's create another feature flag
	// Advance time a bit since versions are based on timestamps.
	clock.Add(time.Second)

	err = suite.stub.CreateFeatureToggle(ctx)
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload2, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload2.Version, payload1.Version)
	require.Equal(suite.T(), payload2.SettingsVersion, payload1.SettingsVersion)
	require.Equal(suite.T(), len(payload2.FeatureToggles), 1)
}

func (suite *TogglesTestSuite) TestListenWithCreateSetting() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{SettingsVersion: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload1, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload1.SettingsVersion, int64(0))

	clock := suite.stub.App.App.Clock().(*clock.Mock)
	// Let's create another dynamic setting
	// Advance time a bit since versions are based on timestamps.
	clock.Add(time.Second)

	err = suite.stub.CreateDynamicSetting(ctx)
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload2, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload2.SettingsVersion, payload1.SettingsVersion)
	require.Equal(suite.T(), payload2.Version, payload1.Version)
	require.Equal(suite.T(), len(payload2.DynamicSettings), 1)
}

func (suite *TogglesTestSuite) TestListenWithUpdate() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1, SettingsVersion: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload1, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload1.Version, int64(0))
	require.Greater(suite.T(), payload1.SettingsVersion, int64(0))

	clock := suite.stub.App.App.Clock().(*clock.Mock)

	// Advance time a bit since versions are based on timestamps.
	clock.Add(time.Second)

	ft := payload1.FeatureToggles[0]
	ft.Description = randomdata.FullName(randomdata.RandomGender)
	err = suite.stub.UpdateFeatureToggle(ctx, ft)
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload2, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Equal(suite.T(), len(payload2.FeatureToggles), 1)
	require.Greater(suite.T(), payload2.Version, payload1.Version)
	require.Equal(suite.T(), payload2.FeatureToggles[0].Description, ft.Description)

	// Update dynamic setting
	ds := payload1.DynamicSettings[0]
	ds.Description = randomdata.FullName(randomdata.RandomGender)
	err = suite.stub.UpdateDynamicSetting(ctx, ds)
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload3, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Equal(suite.T(), len(payload3.DynamicSettings), 1)
	require.Greater(suite.T(), payload3.SettingsVersion, payload1.SettingsVersion)
	require.Equal(suite.T(), payload3.DynamicSettings[0].Description, ds.Description)

}

func (suite *TogglesTestSuite) TestListenWithDelete() {
	ctx := context.Background()
	token, err := suite.stub.App.AuthClient.Authenticate(suite.stub.WithAPiKey(ctx), &pb_auth.AuthenticateRequest{Version: version})
	require.Nil(suite.T(), err, "%+v", err)

	client, err := suite.stub.App.TogglesClient.Listen(suite.stub.WithJwtToken(ctx, token.AccessToken), &pb_toggles.ListenRequest{Version: 1})
	require.Nil(suite.T(), err, "%+v", err)
	payload1, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Greater(suite.T(), payload1.Version, int64(0))

	clock := suite.stub.App.App.Clock().(*clock.Mock)

	// Advance time a bit since versions are based on timestamps.
	clock.Add(time.Second)

	ft := payload1.FeatureToggles[0]
	err = suite.stub.DeleteFeatureToggle(ctx, ids.ID(ft.Id))
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload2, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Equal(suite.T(), 1, len(payload2.FeatureToggles))
	require.Greater(suite.T(), payload2.Version, payload1.Version)
	require.NotNil(suite.T(), payload2.FeatureToggles[0].DeletedAt)
	require.Equal(suite.T(), payload2.FeatureToggles[0].Description, "")

	ds := payload1.DynamicSettings[0]
	err = suite.stub.DeleteDynamicSetting(ctx, ids.ID(ds.Id))
	require.Nil(suite.T(), err, "%+v", err)

	// Advance time and make sure we get the update
	clock.Add(toggles.PollInterval * 2)
	payload3, err := client.Recv()
	require.Nil(suite.T(), err, "%+v", err)
	require.Equal(suite.T(), 1, len(payload3.DynamicSettings))
	require.Greater(suite.T(), payload3.SettingsVersion, payload1.SettingsVersion)
	require.NotNil(suite.T(), payload3.DynamicSettings[0].DeletedAt)
	require.Equal(suite.T(), payload3.DynamicSettings[0].Description, "")

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TogglesTestSuite))
}
