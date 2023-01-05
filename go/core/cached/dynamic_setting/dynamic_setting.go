package cached_dynamic_setting

import (
	"context"
	"fmt"
	"platform/go/core/app"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/dynamic_settings"

	pb_private "platform/go/proto/private"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetEnvironmentVersion(ctx context.Context, id ids.ID, app app.App) (*pb_private.EnvironmentVersion, error) {
	pb, err := app.KV().GetProto(ctx, kv.EnvironmentSettingsVersion, string(id))
	if err != nil {
		// Fetch it from the database
		version, err := dynamic_settings.MaxVersionForEnv(ctx, id, app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot settings for environment")
			}
			return nil, status.Errorf(codes.Internal, "could not get dynamic settings version: %+v", err.Error())
		}
		pb = &pb_private.EnvironmentVersion{
			Id:      string(id),
			Version: version,
		}
		// Populate the cache
		if err := app.KV().SetProto(ctx, kv.EnvironmentSettingsVersion, string(id), pb); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	return pb.(*pb_private.EnvironmentVersion), nil
}

func GetDynamicSettings(ctx context.Context, id ids.ID, app app.App, start, end int64) (*pb_private.EnvironmentDynamicSettings, error) {
	suffix := kv.WithSuffix(fmt.Sprintf("%d-%d", start, end))
	pb, err := app.KV().GetProto(ctx, kv.EnvironmentSettings, string(id), suffix)
	if err != nil {
		// Fetch it from the database
		dsEnvs, err := dynamic_settings.ListForEnv(ctx, id, app.DB(), dynamic_settings.WithStartVersion(start), dynamic_settings.WithEndVersion(end), dynamic_settings.WithDeleted())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot find dynamic settings")
			}
			return nil, status.Errorf(codes.Internal, "invalid settings: %+v", err.Error())
		}
		dss, err := dynamic_settings.MultiPb(ctx, dsEnvs, app.Ory(), dynamic_settings.PbOpts{FillUser: false})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get dynamic settings: %+v", err)
		}

		pb = &pb_private.EnvironmentDynamicSettings{
			Id:              string(id),
			StartingVersion: start,
			EndingVersion:   end,
			DynamicSettings: dss,
		}
		// Populate the cache
		if err := app.KV().SetProto(ctx, kv.EnvironmentSettings, string(id), pb, suffix); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	return pb.(*pb_private.EnvironmentDynamicSettings), nil
}
