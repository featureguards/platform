package stubs

import (
	"context"

	pb_dashboard "platform/go/proto/dashboard"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) createApiKey(ctx context.Context) error {
	name := randomdata.Alphanumeric(8)
	envID := s.Proj.Environments[0].Id
	if _, err := s.App.DashboardClient.CreateApiKey(s.WithToken(ctx), &pb_dashboard.CreateApiKeyRequest{
		Name:          name,
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
