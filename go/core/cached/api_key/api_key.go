package cached_api_key

import (
	"context"
	"platform/go/core/app"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/api_keys"

	pb_project "platform/go/proto/project"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Get(ctx context.Context, id ids.ID, app app.App) (*pb_project.ApiKey, error) {
	pb, err := app.KV().GetProto(ctx, kv.ApiKey, string(id))
	if err != nil {
		// Fetch it from the database
		model, err := api_keys.Get(ctx, id, app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, "cannot find api key")
			}

			return nil, status.Error(codes.Internal, "invalid api key")
		}
		pb, err = api_keys.Pb(model)
		if err != nil {
			return nil, status.Error(codes.Internal, "invalid api key")
		}
		// Populate the cache
		if err := app.KV().SetProto(ctx, kv.ApiKey, string(id), pb); err != nil {
			log.Warningf("%s\n", err)
		}
	}
	return pb.(*pb_project.ApiKey), nil
}
