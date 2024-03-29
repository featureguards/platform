package dashboard

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/environments"
	"platform/go/core/models/feature_toggles"
	"platform/go/core/models/projects"
	"platform/go/core/models/users"
	pb_dashboard "platform/go/proto/dashboard"
	pb_project "platform/go/proto/project"
	"sort"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
	pb_platform "github.com/featureguards/featureguards-go/v2/proto/platform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *DashboardServer) CreateFeatureToggle(ctx context.Context, req *pb_dashboard.CreateFeatureToggleRequest) (*empty.Empty, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}
	if err := s.validate(req.Feature); err != nil {
		return nil, err
	}

	user, err := s.authProject(ctx, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN})
	if err != nil {
		return nil, err
	}

	id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.FeatureToggle)
	if err != nil {
		return nil, err
	}

	// We must release the Redis locks AFTER the transaction is committed/rolledback. Hence, doing
	// it outside the Transaction call below.
	var pendings []*kv.Pending
	defer func() {
		for _, p := range pendings {
			p.Finish(ctx)
		}
	}()
	if err := s.app.DB().Transaction(func(tx *gorm.DB) error {
		proj, err := projects.GetProject(ctx, ids.ID(req.ProjectId), tx, true)
		if err != nil {
			return err
		}
		if len(proj.Environments) <= 0 {
			return status.Error(codes.InvalidArgument, "No environments exist for project")
		}

		if _, err := feature_toggles.GetByName(ctx, proj.ID, req.Feature.Name, tx, feature_toggles.GetFTOpts{}); err == nil || err != models.ErrNotFound {
			if err == nil {
				return status.Error(codes.InvalidArgument, "Feature flag name already exists")
			}
			return status.Errorf(codes.Internal, "Could not create feature flag")
		}

		var isAndroid, isIOS, isWeb, isServer bool
		for _, platform := range req.Feature.Platforms {
			switch platform {
			case pb_platform.Type_ANDROID:
				isAndroid = true
			case pb_platform.Type_IOS:
				isIOS = true
			case pb_platform.Type_WEB:
				isWeb = true
			case pb_platform.Type_DEFAULT:
				isServer = true
			}
		}
		ft := models.FeatureToggle{
			Model:       models.Model{ID: id},
			Name:        req.Feature.Name,
			Description: req.Feature.Description,
			ProjectID:   ids.ID(req.ProjectId),
			CreatedByID: user.ID,
			Type:        req.Feature.ToggleType,
			IsAndroid:   isAndroid,
			IsIOS:       isIOS,
			IsWeb:       isWeb,
			IsServer:    isServer,
		}

		var ftEnvs []models.FeatureToggleEnv
		// Sort the keys to avoid dead-locking.
		sort.Slice(proj.Environments, func(i, j int) bool {
			return proj.Environments[i].ID < proj.Environments[j].ID
		})
		for _, env := range proj.Environments {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentVersion, env.ID.String())
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)
			ftEnvID, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.FeatureToggleEnv)
			if err != nil {
				return err
			}
			proto, err := feature_toggles.SerializeDefinition(ctx, req.Feature)
			if err != nil {
				return err
			}
			ftEnv := models.FeatureToggleEnv{
				Model:           models.Model{ID: ftEnvID},
				EnvironmentID:   env.ID,
				ProjectID:       ids.ID(req.ProjectId),
				FeatureToggleID: id,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     s.app.Clock().Now().UnixNano(),
				Enabled:     req.Feature.Enabled,
				CreatedByID: user.ID,
				Proto:       proto,
			}
			ftEnvs = append(ftEnvs, ftEnv)
		}
		if err := tx.Create(&ft).Error; err != nil {
			return err
		}

		if err := tx.Create(&ftEnvs).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "could not create feature flag")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) ListFeatureToggles(ctx context.Context, req *pb_dashboard.ListFeatureToggleRequest) (*pb_dashboard.ListFeatureToggleResponse, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	if _, err := s.authEnvironment(ctx, ids.ID(req.EnvironmentId)); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not list feature flags")
	}

	// Must pass since permissions are based on the project ID.
	ftEnvs, err := feature_toggles.ListForEnv(ctx, ids.ID(req.EnvironmentId), s.app.DB())
	if err != nil {
		if err == models.ErrNotFound {
			return &pb_dashboard.ListFeatureToggleResponse{FeatureToggles: make([]*pb_ft.FeatureToggle, 0)}, nil
		}
		return nil, status.Errorf(codes.Internal, "could not get feature flag")
	}

	fts, err := feature_toggles.MultiPb(ctx, ftEnvs, s.app.Ory(), feature_toggles.PbOpts{FillUser: false})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get feature flag")
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
		envs, err := environments.List(ctx, ft.ProjectID, s.app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Error(codes.NotFound, "no feature flag found")
			}
			return nil, status.Errorf(codes.Internal, "could not get feature flag")
		}
		for _, env := range envs {
			envIDs = append(envIDs, env.ID)
		}
	}
	if len(envIDs) <= 0 {
		return nil, status.Error(codes.NotFound, "no feature flag found")
	}

	res := make([]*pb_dashboard.EnvironmentFeatureToggle, 0, len(envIDs))
	for _, envID := range envIDs {
		// There is no easy query to make this because there could be an imbalance of versions
		// across environments
		ftEnv, err := feature_toggles.GetLatestForEnv(ctx, ids.ID(req.Id), envID, s.app.DB())
		if err != nil {
			if err == models.ErrNotFound {
				return nil, status.Errorf(codes.NotFound, "feature flag not found")
			}
			return nil, status.Errorf(codes.Internal, "could not get feature flag")
		}
		pb, err := feature_toggles.Pb(ctx, ftEnv, s.app.Ory(), feature_toggles.PbOpts{FillUser: true})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature flag")
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

	ftEnvs, err := feature_toggles.GetHistoryForEnv(ctx, ids.ID(req.Id), ids.ID(req.EnvironmentId), s.app.DB())
	if err != nil {
		if err == models.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "feature flag not found")
		}
		return nil, status.Errorf(codes.Internal, "could not get feature flag history")
	}
	var fts []*pb_ft.FeatureToggle
	for _, ftEnv := range ftEnvs {
		pb, err := feature_toggles.Pb(ctx, &ftEnv, s.app.Ory(), feature_toggles.PbOpts{FillUser: true})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get feature flag")
		}
		fts = append(fts, pb)
	}
	return &pb_ft.FeatureToggleHistory{
		History: fts,
	}, nil
}

func (s *DashboardServer) UpdateFeatureToggle(ctx context.Context, req *pb_dashboard.UpdateFeatureToggleRequest) (*empty.Empty, error) {
	if req.Feature == nil || req.Feature.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}
	if req.Feature.Id != req.Id {
		return nil, status.Error(codes.InvalidArgument, "IDs must match")
	}
	existing, err := s.authFeatureToggle(ctx, ids.ID(req.Feature.Id))
	if err != nil {
		return nil, err
	}
	if existing.ID != ids.ID(req.Feature.Id) || existing.Name != req.Feature.Name || existing.Type != req.Feature.ToggleType || existing.ProjectID != ids.ID(req.Feature.ProjectId) {
		return nil, status.Error(codes.InvalidArgument, "ID, name, project_id and toggle_type cannot be changed")
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB())
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	if err := s.validate(req.Feature); err != nil {
		return nil, err
	}

	// We must release the Redis locks AFTER the transaction is committed/rolledback. Hence, doing
	// it outside the Transaction call below.
	var pendings []*kv.Pending
	defer func() {
		for _, p := range pendings {
			p.Finish(ctx)
		}
	}()
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		proj, err := projects.GetProject(ctx, existing.ProjectID, tx, true)
		if err != nil {
			return err
		}
		if len(proj.Environments) <= 0 {
			return status.Error(codes.InvalidArgument, "no environments exist for project")
		}
		// Make sure environment IDs match
		projectEnvIDs := make(map[ids.ID]struct{}, len(proj.Environments))
		for _, env := range proj.Environments {
			projectEnvIDs[env.ID] = struct{}{}
		}
		for _, envID := range req.EnvironmentIds {
			if _, ok := projectEnvIDs[ids.ID(envID)]; !ok {
				return status.Error(codes.InvalidArgument, "environment not found")
			}
		}

		existing.Description = req.Feature.Description
		existing.IsAndroid = false
		existing.IsIOS = false
		existing.IsWeb = false
		for _, platform := range req.Feature.Platforms {
			switch platform {
			case pb_platform.Type_ANDROID:
				existing.IsAndroid = true
			case pb_platform.Type_IOS:
				existing.IsIOS = true
			case pb_platform.Type_WEB:
				existing.IsWeb = true
			}
		}

		var ftEnvs []models.FeatureToggleEnv
		// Sort the keys to avoid dead-locking.
		sort.Strings(req.EnvironmentIds)
		for _, envID := range req.EnvironmentIds {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentVersion, envID)
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)
			ftEnvID, err := ids.IDFromRoot(existing.ProjectID, ids.FeatureToggleEnv)
			if err != nil {
				return err
			}
			proto, err := feature_toggles.SerializeDefinition(ctx, req.Feature)
			if err != nil {
				return err
			}
			ftEnv := models.FeatureToggleEnv{
				Model:           models.Model{ID: ftEnvID},
				EnvironmentID:   ids.ID(envID),
				ProjectID:       existing.ProjectID,
				FeatureToggleID: existing.ID,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     s.app.Clock().Now().UnixNano(),
				Enabled:     req.Feature.Enabled,
				CreatedByID: user.ID,
				Proto:       proto,
			}
			ftEnvs = append(ftEnvs, ftEnv)
		}
		if err := tx.Save(&existing).Error; err != nil {
			return errors.WithStack(err)
		}

		if err := tx.Create(&ftEnvs).Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Errorf(codes.Internal, "could not update feature flag")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) DeleteFeatureToggle(ctx context.Context, req *pb_dashboard.DeleteFeatureToggleRequest) (*empty.Empty, error) {
	ft, err := s.authFeatureToggle(ctx, ids.ID(req.Id))
	if err != nil {
		return nil, err
	}
	// We must release the Redis locks AFTER the transaction is committed/rolledback. Hence, doing
	// it outside the Transaction call below.
	var pendings []*kv.Pending
	defer func() {
		for _, p := range pendings {
			p.Finish(ctx)
		}
	}()
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// Take a lock on the project
		proj, err := projects.GetProject(ctx, ft.ProjectID, tx, true)
		if err != nil {
			return err
		}
		// Sort the keys to avoid dead-locking.
		sort.Slice(proj.Environments, func(i, j int) bool {
			return proj.Environments[i].ID < proj.Environments[j].ID
		})
		for _, env := range proj.Environments {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentVersion, env.ID.String())
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)
		}
		// Must pass since permissions are based on the project ID.
		if err := tx.Delete(&models.FeatureToggle{
			Model:     models.Model{ID: ids.ID(req.Id)},
			ProjectID: ft.ProjectID,
		}).Error; err != nil {
			return errors.WithStack(err)
		}
		// Delete it from all environments and all versions.
		if err := tx.Model(&models.FeatureToggleEnv{}).Where(
			"feature_toggle_id = ? AND project_id = ?", req.Id, ft.ProjectID).Update("version", s.app.Clock().Now().UnixNano()).Error; err != nil {
			return errors.WithStack(err)
		}
		if err := tx.Delete(&models.FeatureToggleEnv{}, "feature_toggle_id = ? AND project_id = ?", req.Id, ft.ProjectID).Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Error(codes.Internal, "could not delete feature flag")

	}
	return &empty.Empty{}, nil
}

func (s *DashboardServer) validate(ft *pb_ft.FeatureToggle) error {
	if ft.Name == "" {
		return status.Error(codes.InvalidArgument, "Name is not specified")
	}
	if len(ft.Platforms) == 0 {
		return status.Error(codes.InvalidArgument, "No platform specified")
	}

	switch ft.ToggleType {
	case pb_ft.FeatureToggle_PERCENTAGE:
		pcnt := ft.GetPercentage()
		if pcnt == nil {
			return status.Error(codes.InvalidArgument, "No percentage definition")
		}
		if pcnt.Stickiness.StickinessType == pb_ft.Stickiness_KEYS && len(pcnt.Stickiness.Keys) == 0 {
			return status.Error(codes.InvalidArgument, "No attributes specified for a sticky percentage feature flag.")
		}
	}
	return nil
}
