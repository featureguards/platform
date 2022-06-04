package stubs

import (
	"context"

	"platform/go/core/ids"
	pb_dashboard "platform/go/proto/dashboard"

	pb_ft "github.com/featureguards/client-go/proto/feature_toggle"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) CreateFeatureToggle(ctx context.Context) error {
	name := randomdata.Noun()
	_, err := s.App.DashboardClient.CreateFeatureToggle(s.WithToken(ctx), &pb_dashboard.CreateFeatureToggleRequest{
		Feature: &pb_ft.FeatureToggle{
			Name:       name,
			ToggleType: pb_ft.FeatureToggle_ON_OFF,
			Enabled:    true,
		},
		ProjectId: s.Proj.Id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Stubs) UpdateFeatureToggle(ctx context.Context, ft *pb_ft.FeatureToggle) error {
	var envIDs []string
	for _, env := range s.Proj.Environments {
		envIDs = append(envIDs, env.Id)
	}
	_, err := s.App.DashboardClient.UpdateFeatureToggle(s.WithToken(ctx), &pb_dashboard.UpdateFeatureToggleRequest{
		Id:             ft.Id,
		Feature:        ft,
		EnvironmentIds: envIDs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Stubs) DeleteFeatureToggle(ctx context.Context, id ids.ID) error {
	_, err := s.App.DashboardClient.DeleteFeatureToggle(s.WithToken(ctx), &pb_dashboard.DeleteFeatureToggleRequest{
		Id: string(id),
	})
	if err != nil {
		return err
	}

	return nil
}
