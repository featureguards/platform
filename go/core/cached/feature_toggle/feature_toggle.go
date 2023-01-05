package cached_feature_toggle

import (
	"context"
	"fmt"
	"platform/go/core/app"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/feature_toggles"

	pb_private "platform/go/proto/private"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetEnvironmentVersion(ctx context.Context, id ids.ID, app app.App) (*pb_private.EnvironmentVersion, error) {
	pb, err := app.KV().GetProto(ctx, kv.EnvironmentVersion, string(id))
	if err != nil {
		// Fetch it from the database
		version, err := feature_toggles.MaxVersionForEnv(ctx, id, app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot feature toggles for environment")
			}
			return nil, status.Error(codes.Internal, "invalid feature toggle version")
		}
		pb = &pb_private.EnvironmentVersion{
			Id:      string(id),
			Version: version,
		}
		// Populate the cache
		if err := app.KV().SetProto(ctx, kv.EnvironmentVersion, string(id), pb); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	return pb.(*pb_private.EnvironmentVersion), nil
}

func GetFeatureToggles(ctx context.Context, id ids.ID, app app.App, start, end int64) (*pb_private.EnvironmentFeatureToggles, error) {
	suffix := kv.WithSuffix(fmt.Sprintf("%d-%d", start, end))
	pb, err := app.KV().GetProto(ctx, kv.EnvironmentToggles, string(id), suffix)
	if err != nil {
		// Fetch it from the database
		ftEnvs, err := feature_toggles.ListForEnv(ctx, id, app.DB(), feature_toggles.WithStartVersion(start), feature_toggles.WithEndVersion(end), feature_toggles.WithDeleted())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot feature flags")
			}
			return nil, status.Error(codes.Internal, "invalid feature flags")
		}
		fts, err := feature_toggles.MultiPb(ctx, ftEnvs, app.Ory(), feature_toggles.PbOpts{FillUser: false})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature flag")
		}

		pb = &pb_private.EnvironmentFeatureToggles{
			Id:              string(id),
			StartingVersion: start,
			EndingVersion:   end,
			FeatureToggles:  fts,
		}
		// Populate the cache
		if err := app.KV().SetProto(ctx, kv.EnvironmentToggles, string(id), pb, suffix); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	return pb.(*pb_private.EnvironmentFeatureToggles), nil
}
