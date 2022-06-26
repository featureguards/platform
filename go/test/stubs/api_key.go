package stubs

import (
	"context"

	pb_dashboard "platform/go/proto/dashboard"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) createApiKey(ctx context.Context) error {
	name := randomdata.Noun()
	envID := s.Proj.Environments[0].Id
	if _, err := s.App.DashboardClient.CreateApiKey(s.WithToken(ctx), &pb_dashboard.CreateApiKeyRequest{
		Name:          name,
		Platforms:     []pb_ft.Platform_Type{pb_ft.Platform_WEB},
		EnvironmentId: envID,
	}); err != nil {
		return err
	}
	keys, err := s.App.DashboardClient.ListApiKeys(s.WithToken(ctx), &pb_dashboard.ListApiKeysRequest{EnvironmentId: envID})
	if err != nil {
		return err
	}

	s.ApiKey = keys.ApiKeys[0]
	return nil
}
