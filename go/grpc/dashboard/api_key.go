package dashboard

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/api_keys"
	"platform/go/core/models/projects"
	"platform/go/core/random"
	pb_dashboard "platform/go/proto/dashboard"
	pb_project "platform/go/proto/project"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *DashboardServer) CreateApiKey(ctx context.Context, req *pb_dashboard.CreateApiKeyRequest) (*emptypb.Empty, error) {
	// We validate here
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	env, err := s.authEnvironment(ctx, ids.ID(req.EnvironmentId))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not find environment")
	}

	id, err := ids.IDFromRoot(ids.ID(env.ProjectID), ids.ApiKey)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "could not create environment")
	}

	m := models.ApiKey{
		Model:         models.Model{ID: id},
		Name:          req.Name,
		ProjectID:     ids.ID(env.ProjectID),
		EnvironmentID: ids.ID(env.ID),
		Key:           id.String() + ":" + random.RandString(16, nil),
	}

	if req.ExpiresAt != nil {
		m.ExpiresAt = req.ExpiresAt.AsTime()
	}

	// Take a lock
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := projects.GetProject(ctx, ids.ID(env.ProjectID), tx, true); err != nil {
			return err
		}
		if err := tx.Create(&m).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Error(codes.Internal, "could not create api key")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (s *DashboardServer) ListApiKeys(ctx context.Context, req *pb_dashboard.ListApiKeysRequest) (*pb_project.ApiKeys, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}

	if _, err := s.authEnvironment(ctx, ids.ID(req.EnvironmentId)); err != nil {
		return nil, err
	}

	apiKeys, err := api_keys.List(ctx, ids.ID(req.EnvironmentId), s.DB(ctx))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no environments found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not list api keys")
	}
	var pbApiKeys []*pb_project.ApiKey
	for _, apiKey := range apiKeys {
		pbApiKey, err := api_keys.Pb(&apiKey)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not list environments")
		}
		pbApiKeys = append(pbApiKeys, pbApiKey)
	}
	return &pb_project.ApiKeys{
		ApiKeys: pbApiKeys,
	}, nil
}

func (s *DashboardServer) DeleteApiKey(ctx context.Context, req *pb_dashboard.DeleteApiKeyRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}
	apiKey, err := api_keys.Get(ctx, ids.ID(req.Id), s.app.DB())
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no api key found")
		}
		return nil, status.Error(codes.Internal, "could not delete api key")
	}

	if _, err := s.authProject(ctx, apiKey.ProjectID, allRoles); err != nil {
		return nil, err
	}
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the project
		_, err := projects.GetProject(ctx, apiKey.ProjectID, tx, true)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return status.Error(codes.NotFound, "no projec found")
			}
			return status.Error(codes.Internal, "could not delete api key")
		}

		pending, err := s.app.KV().StartPending(ctx, kv.ApiKey, req.Id)
		if err != nil {
			log.Errorf("%s\n", err)
			return status.Error(codes.Internal, "could not delete api key")
		}
		defer pending.Finish(ctx)
		if err := tx.Delete(&models.ApiKey{Model: models.Model{ID: ids.ID(req.Id)}}).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Error(codes.Internal, "could not delete api key")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}
