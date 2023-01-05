package stubs

import (
	"context"

	"platform/go/core/ids"
	pb_dashboard "platform/go/proto/dashboard"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) CreateDynamicSetting(ctx context.Context) error {
	name := randomdata.Noun()
	_, err := s.App.DashboardClient.CreateDynamicSetting(s.WithToken(ctx), &pb_dashboard.CreateDynamicSettingRequest{
		Setting: &pb_ds.DynamicSetting{
			Name:        name,
			SettingType: pb_ds.DynamicSetting_BOOL,
			Platforms:   s.ApiKey.Platforms,
			SettingDefinition: &pb_ds.DynamicSetting_BoolValue{
				BoolValue: &pb_ds.BoolValue{Value: true}},
		},
		ProjectId: s.Proj.Id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Stubs) UpdateDynamicSetting(ctx context.Context, ds *pb_ds.DynamicSetting) error {
	var envIDs []string
	for _, env := range s.Proj.Environments {
		envIDs = append(envIDs, env.Id)
	}
	_, err := s.App.DashboardClient.UpdateDynamicSetting(s.WithToken(ctx), &pb_dashboard.UpdateDynamicSettingRequest{
		Id:             ds.Id,
		Setting:        ds,
		EnvironmentIds: envIDs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Stubs) DeleteDynamicSetting(ctx context.Context, id ids.ID) error {
	_, err := s.App.DashboardClient.DeleteDynamicSetting(s.WithToken(ctx), &pb_dashboard.DeleteDynamicSettingRequest{
		Id: string(id),
	})
	if err != nil {
		return err
	}

	return nil
}
