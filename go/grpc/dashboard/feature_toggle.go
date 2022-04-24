package dashboard

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/environments"
	"stackv2/go/core/models/feature_toggles"
	"stackv2/go/core/models/projects"
	"stackv2/go/core/models/users"
	pb_dashboard "stackv2/go/proto/dashboard"
	pb_ft "stackv2/go/proto/feature_toggle"
	pb_project "stackv2/go/proto/project"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *DashboardServer) CreateFeatureToggle(ctx context.Context, req *pb_dashboard.CreateFeatureToggleRequest) (*empty.Empty, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}
	if req.Feature.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is not specified")
	}

	user, err := s.authProject(ctx, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN})
	if err != nil {
		return nil, err
	}

	id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.FeatureToggle)
	if err != nil {
		return nil, err
	}
	proj, err := projects.GetProject(ctx, ids.ID(req.ProjectId), s.app.DB)
	if err != nil {
		return nil, err
	}
	if len(proj.Environments) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no environments exist for project")
	}

	if _, err := feature_toggles.GetByName(ctx, proj.ID, req.Feature.Name, s.app.DB, feature_toggles.GetFTOpts{}); err == nil || err != models.ErrNotFound {
		if err == nil {
			return nil, status.Error(codes.InvalidArgument, "feature toggle name already exists")
		}
		return nil, status.Errorf(codes.Internal, "could not create feature toggle")
	}

	if err := s.app.DB.Transaction(func(tx *gorm.DB) error {
		var isMobile, isWeb bool
		for _, platform := range req.Feature.Platforms {
			switch platform {
			case pb_ft.Platform_MOBILE:
				isMobile = true
				break
			case pb_ft.Platform_WEB:
				isWeb = true
				break
			}
		}
		ft := models.FeatureToggle{
			Model:       models.Model{ID: id},
			Name:        req.Feature.Name,
			Description: req.Feature.Description,
			ProjectID:   ids.ID(req.ProjectId),
			CreatedByID: user.ID,
			Type:        req.Feature.ToggleType,
			IsMobile:    isMobile,
			IsWeb:       isWeb,
		}

		var ftEnvs []models.FeatureToggleEnv
		for _, env := range proj.Environments {
			ftEnvID, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.FeatureToggleEnv)
			if err != nil {
				return err
			}
			proto, err := feature_toggles.SerializeDefinition(ctx, req.Feature)
			if err != nil {
				log.Error(errors.WithStack(err))
				return err
			}
			ftEnv := models.FeatureToggleEnv{
				Model:           models.Model{ID: ftEnvID},
				EnvironmentID:   env.ID,
				ProjectID:       ids.ID(req.ProjectId),
				FeatureToggleID: id,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     time.Now().UnixNano(),
				Enabled:     req.Feature.Enabled,
				CreatedByID: user.ID,
				Proto:       proto,
			}
			ftEnvs = append(ftEnvs, ftEnv)
		}
		if err := tx.Create(&ft).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}

		if err := tx.Create(&ftEnvs).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}
		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "could not create feature toggle")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) ListFeatureToggles(ctx context.Context, req *pb_dashboard.ListFeatureToggleRequest) (*pb_dashboard.ListFeatureToggleResponse, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	env, err := environments.Get(ctx, ids.ID(req.EnvironmentId), s.app.DB)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not list feature toggles")
	}

	if _, err := s.authProject(ctx, ids.ID(env.ProjectID), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}
	// Must pass since permissions are based on the project ID.
	ftEnvs, err := feature_toggles.ListLatestForEnv(ctx, ids.ID(req.EnvironmentId), s.app.DB)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "feature toggle not found")
		}
		return nil, status.Errorf(codes.Internal, "could not get feature toggle")
	}

	var fts []*pb_ft.FeatureToggle
	for _, ftEnv := range ftEnvs {
		pb, err := feature_toggles.Pb(ctx, &ftEnv, s.app.Ory, feature_toggles.PbOpts{FillUser: false})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature toggle")
		}
		fts = append(fts, pb)
	}
	return &pb_dashboard.ListFeatureToggleResponse{
		FeatureToggles: fts,
	}, nil
}
func (s *DashboardServer) GetFeatureToggle(ctx context.Context, req *pb_dashboard.GetFeatureToggleRequest) (*pb_dashboard.EnvironmentFeatureToggles, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	ft, err := s.authFeatureToggle(ctx, ids.ID(req.Id))
	if err != nil {
		return nil, err
	}

	envIDs := ids.FromStringSlice(req.EnvironmentIds)
	if len(envIDs) <= 0 {
		// We will feat all for the project.
		envs, err := environments.List(ctx, ft.ProjectID, s.app.DB)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Error(codes.NotFound, "no feature toggle found")
			}
			return nil, status.Errorf(codes.Internal, "could not get feature toggle")
		}
		for _, env := range envs {
			envIDs = append(envIDs, env.ID)
		}
	}
	if len(envIDs) <= 0 {
		return nil, status.Error(codes.NotFound, "no feature toggle found")
	}

	res := make([]*pb_dashboard.EnvironmentFeatureToggle, 0, len(envIDs))
	for _, envID := range envIDs {
		// There is no easy query to make this because there could be an imbalance of versions
		// across environments
		ftEnv, err := feature_toggles.GetLatestForEnv(ctx, ids.ID(req.Id), envID, s.app.DB)
		if err != nil {
			if err == models.ErrNotFound {
				return nil, status.Errorf(codes.NotFound, "feature toggle not found")
			}
			return nil, status.Errorf(codes.Internal, "could not get feature toggle")
		}
		pb, err := feature_toggles.Pb(ctx, ftEnv, s.app.Ory, feature_toggles.PbOpts{FillUser: true})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature toggle")
		}
		res = append(res, &pb_dashboard.EnvironmentFeatureToggle{EnvironmentId: string(envID), FeatureToggle: pb})
	}

	return &pb_dashboard.EnvironmentFeatureToggles{FeatureToggles: res}, nil
}

func (s *DashboardServer) GetFeatureToggleHistoryForEnvironment(ctx context.Context, req *pb_dashboard.GetFeatureToggleHistoryRequest) (*pb_ft.FeatureToggleHistory, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	if _, err := s.authFeatureToggle(ctx, ids.ID(req.Id)); err != nil {
		return nil, err
	}

	ftEnvs, err := feature_toggles.GetHistoryForEnv(ctx, ids.ID(req.Id), ids.ID(req.EnvironmentId), s.app.DB)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "feature toggle not found")
		}
		return nil, status.Errorf(codes.Internal, "could not get feature toggle history")
	}
	var fts []*pb_ft.FeatureToggle
	for _, ftEnv := range ftEnvs {
		pb, err := feature_toggles.Pb(ctx, &ftEnv, s.app.Ory, feature_toggles.PbOpts{FillUser: true})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature toggle")
		}
		fts = append(fts, pb)
	}
	return &pb_ft.FeatureToggleHistory{
		History: fts,
	}, nil
}

func (s *DashboardServer) UpdateFeatureToggle(ctx context.Context, req *pb_dashboard.UpdateFeatureToggleRequest) (*empty.Empty, error) {
	if req.Feature == nil || req.Feature.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}
	if req.Feature.Id != req.Id {
		return nil, status.Error(codes.InvalidArgument, "ids must match")
	}
	existing, err := s.authFeatureToggle(ctx, ids.ID(req.Feature.Id))
	if err != nil {
		return nil, err
	}
	if existing.ID != ids.ID(req.Feature.Id) || existing.Name != req.Feature.Name || existing.Type != req.Feature.ToggleType || existing.ProjectID != ids.ID(req.Feature.ProjectId) {
		return nil, status.Error(codes.InvalidArgument, "id, name, project_id and toggle_type cannot be changed")
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	proj, err := projects.GetProject(ctx, existing.ProjectID, s.app.DB)
	if err != nil {
		return nil, err
	}
	if len(proj.Environments) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no environments exist for project")
	}
	// Make sure environment IDs match
	projectEnvIDs := make(map[ids.ID]struct{}, len(proj.Environments))
	for _, env := range proj.Environments {
		projectEnvIDs[env.ID] = struct{}{}
	}
	for _, envID := range req.EnvironmentIds {
		if _, ok := projectEnvIDs[ids.ID(envID)]; !ok {
			return nil, status.Error(codes.InvalidArgument, "environment not found")
		}
	}

	// TODO: Add validation for definition
	if err := s.app.DB.Transaction(func(tx *gorm.DB) error {
		existing.Description = req.Feature.Description
		existing.IsMobile = false
		existing.IsWeb = false
		for _, platform := range req.Feature.Platforms {
			switch platform {
			case pb_ft.Platform_MOBILE:
				existing.IsMobile = true
				break
			case pb_ft.Platform_WEB:
				existing.IsWeb = true
				break
			}
		}

		var ftEnvs []models.FeatureToggleEnv
		for _, envID := range req.EnvironmentIds {
			ftEnvID, err := ids.IDFromRoot(existing.ProjectID, ids.FeatureToggleEnv)
			if err != nil {
				return err
			}
			proto, err := feature_toggles.SerializeDefinition(ctx, req.Feature)
			if err != nil {
				log.Error(errors.WithStack(err))
				return err
			}
			ftEnv := models.FeatureToggleEnv{
				Model:           models.Model{ID: ftEnvID},
				EnvironmentID:   ids.ID(envID),
				ProjectID:       existing.ProjectID,
				FeatureToggleID: existing.ID,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     time.Now().UnixNano(),
				Enabled:     req.Feature.Enabled,
				CreatedByID: user.ID,
				Proto:       proto,
			}
			ftEnvs = append(ftEnvs, ftEnv)
		}
		if err := tx.Save(&existing).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}

		if err := tx.Create(&ftEnvs).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}
		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "could not update feature toggle")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) DeleteFeatureToggle(ctx context.Context, req *pb_dashboard.DeleteFeatureToggleRequest) (*empty.Empty, error) {
	ft, err := s.authFeatureToggle(ctx, ids.ID(req.Id))
	if err != nil {
		return nil, err
	}
	if err := s.app.DB.Transaction(func(tx *gorm.DB) error {
		// Must pass since permissions are based on the project ID.
		if err := s.app.DB.WithContext(ctx).Delete(&models.FeatureToggle{
			Model:     models.Model{ID: ids.ID(req.Id)},
			ProjectID: ft.ProjectID,
		}).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}
		// Delete it from all environments and all versions.
		if err := s.app.DB.WithContext(ctx).Delete(&models.FeatureToggleEnv{
			FeatureToggleID: ids.ID(req.Id),
			ProjectID:       ft.ProjectID,
		}).Error; err != nil {
			log.Error(errors.WithStack(err))
			return err
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, "could not delete feature toggle")

	}
	return &empty.Empty{}, nil
}
