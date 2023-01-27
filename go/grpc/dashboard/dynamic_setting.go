package dashboard

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/models"
	"platform/go/core/models/dynamic_settings"
	"platform/go/core/models/environments"
	"platform/go/core/models/projects"
	"platform/go/core/models/users"
	pb_dashboard "platform/go/proto/dashboard"
	pb_project "platform/go/proto/project"
	"sort"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"
	pb_platform "github.com/featureguards/featureguards-go/v2/proto/platform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *DashboardServer) CreateDynamicSetting(ctx context.Context, req *pb_dashboard.CreateDynamicSettingRequest) (*empty.Empty, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}
	if err := s.validateDynamicSetting(req.Setting); err != nil {
		return nil, err
	}

	user, err := s.authProject(ctx, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN})
	if err != nil {
		return nil, err
	}

	id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.DynamicSetting)
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

		if _, err := dynamic_settings.GetByName(ctx, proj.ID, req.Setting.Name, tx, dynamic_settings.GetDSOpts{}); err == nil || err != models.ErrNotFound {
			if err == nil {
				return status.Error(codes.InvalidArgument, "Dynamic setting name already exists")
			}
			return status.Errorf(codes.Internal, "Could not create dynamic setting")
		}

		var isAndroid, isIOS, isWeb, isServer bool
		for _, platform := range req.Setting.Platforms {
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
		ds := models.DynamicSetting{
			Model:       models.Model{ID: id},
			Name:        req.Setting.Name,
			Description: req.Setting.Description,
			ProjectID:   ids.ID(req.ProjectId),
			CreatedByID: user.ID,
			Type:        req.Setting.SettingType,
			IsAndroid:   isAndroid,
			IsIOS:       isIOS,
			IsWeb:       isWeb,
			IsServer:    isServer,
		}

		var dsEnvs []models.DynamicSettingEnv
		// Sort the keys to avoid dead-locking.
		sort.Slice(proj.Environments, func(i, j int) bool {
			return proj.Environments[i].ID < proj.Environments[j].ID
		})
		for _, env := range proj.Environments {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentSettingsVersion, env.ID.String())
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)
			dsEnvID, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.DynamicSettingEnv)
			if err != nil {
				return err
			}
			proto, err := dynamic_settings.SerializeDefinition(ctx, req.Setting)
			if err != nil {
				return err
			}
			dsEnv := models.DynamicSettingEnv{
				Model:            models.Model{ID: dsEnvID},
				EnvironmentID:    env.ID,
				ProjectID:        ids.ID(req.ProjectId),
				DynamicSettingID: id,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     s.app.Clock().Now().UnixNano(),
				CreatedByID: user.ID,
				Proto:       proto,
			}
			dsEnvs = append(dsEnvs, dsEnv)
		}
		if err := tx.Create(&ds).Error; err != nil {
			return err
		}

		if err := tx.Create(&dsEnvs).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "could not create dynamic setting")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) ListDynamicSettings(ctx context.Context, req *pb_dashboard.ListDynamicSettingRequest) (*pb_dashboard.ListDynamicSettingResponse, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	if _, err := s.authEnvironment(ctx, ids.ID(req.EnvironmentId)); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		log.Errorf("%s\n", errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not list dynamic settings")
	}

	// Must pass since permissions are based on the project ID.
	dsEnvs, err := dynamic_settings.ListForEnv(ctx, ids.ID(req.EnvironmentId), s.app.DB())
	if err != nil {
		if err == models.ErrNotFound {
			return &pb_dashboard.ListDynamicSettingResponse{DynamicSettings: make([]*pb_ds.DynamicSetting, 0)}, nil
		}
		log.Errorf("%s\n", errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not get dynamic setting")
	}

	dss, err := dynamic_settings.MultiPb(ctx, dsEnvs, s.app.Ory(), dynamic_settings.PbOpts{FillUser: false})
	if err != nil {
		log.Errorf("%s\n", errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not get dynamic setting")
	}
	return &pb_dashboard.ListDynamicSettingResponse{
		DynamicSettings: dss,
	}, nil
}
func (s *DashboardServer) GetDynamicSetting(ctx context.Context, req *pb_dashboard.GetDynamicSettingRequest) (*pb_dashboard.EnvironmentDynamicSettings, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	ds, err := s.authDynamicSetting(ctx, ids.ID(req.Id))
	if err != nil {
		return nil, err
	}

	envIDs := ids.FromStringSlice(req.EnvironmentIds)
	if len(envIDs) <= 0 {
		// We will feat all for the project.
		envs, err := environments.List(ctx, ds.ProjectID, s.app.DB())
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, status.Error(codes.NotFound, "no dynamic setting found")
			}
			log.Errorf("%s\n", errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not get dynamic setting")
		}
		for _, env := range envs {
			envIDs = append(envIDs, env.ID)
		}
	}
	if len(envIDs) <= 0 {
		return nil, status.Error(codes.NotFound, "no dynamic setting found")
	}

	res := make([]*pb_dashboard.EnvironmentDynamicSetting, 0, len(envIDs))
	for _, envID := range envIDs {
		// There is no easy query to make this because there could be an imbalance of versions
		// across environments
		dsEnv, err := dynamic_settings.GetLatestForEnv(ctx, ids.ID(req.Id), envID, s.app.DB())
		if err != nil {
			if err == models.ErrNotFound {
				return nil, status.Errorf(codes.NotFound, "dynamic setting not found")
			}
			log.Errorf("%s\n", errors.WithStack(err))
			return nil, status.Errorf(codes.Internal, "could not get dynamic setting")
		}
		pb, err := dynamic_settings.Pb(ctx, dsEnv, s.app.Ory(), dynamic_settings.PbOpts{FillUser: true})
		if err != nil {
			log.Errorf("%s\n", errors.WithStack(err))
			return nil, status.Errorf(codes.Internal, "could not get dynamic setting")
		}
		res = append(res, &pb_dashboard.EnvironmentDynamicSetting{EnvironmentId: string(envID), Setting: pb})
	}

	return &pb_dashboard.EnvironmentDynamicSettings{Settings: res}, nil
}

func (s *DashboardServer) GetDynamicSettingHistoryForEnvironment(ctx context.Context, req *pb_dashboard.GetDynamicSettingHistoryRequest) (*pb_ds.DynamicSettingHistory, error) {
	if req.EnvironmentId == "" {
		return nil, status.Error(codes.InvalidArgument, "environment_id is not specified")
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	if _, err := s.authDynamicSetting(ctx, ids.ID(req.Id)); err != nil {
		return nil, err
	}

	dsEnvs, err := dynamic_settings.GetHistoryForEnv(ctx, ids.ID(req.Id), ids.ID(req.EnvironmentId), s.app.DB())
	if err != nil {
		if err == models.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "dynamic setting not found")
		}
		log.Errorf("%s\n", errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not get dynamic setting history")
	}
	var dss []*pb_ds.DynamicSetting
	for _, dsEnv := range dsEnvs {
		pb, err := dynamic_settings.Pb(ctx, &dsEnv, s.app.Ory(), dynamic_settings.PbOpts{FillUser: true})
		if err != nil {
			log.Errorf("%s\n", errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not get dynamic setting")
		}
		dss = append(dss, pb)
	}
	return &pb_ds.DynamicSettingHistory{
		History: dss,
	}, nil
}

func (s *DashboardServer) UpdateDynamicSetting(ctx context.Context, req *pb_dashboard.UpdateDynamicSettingRequest) (*empty.Empty, error) {
	if req.Setting == nil || req.Setting.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is not specified")
	}
	if req.Setting.Id != req.Id {
		return nil, status.Error(codes.InvalidArgument, "IDs must match")
	}
	existing, err := s.authDynamicSetting(ctx, ids.ID(req.Setting.Id))
	if err != nil {
		return nil, err
	}
	if existing.ID != ids.ID(req.Setting.Id) || existing.Name != req.Setting.Name || existing.Type != req.Setting.SettingType || existing.ProjectID != ids.ID(req.Setting.ProjectId) {
		return nil, status.Error(codes.InvalidArgument, "ID, name, project_id and toggle_type cannot be changed")
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB())
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	if err := s.validateDynamicSetting(req.Setting); err != nil {
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

		existing.Description = req.Setting.Description
		existing.IsAndroid = false
		existing.IsIOS = false
		existing.IsWeb = false
		for _, platform := range req.Setting.Platforms {
			switch platform {
			case pb_platform.Type_ANDROID:
				existing.IsAndroid = true
			case pb_platform.Type_IOS:
				existing.IsIOS = true
			case pb_platform.Type_WEB:
				existing.IsWeb = true
			}
		}

		var dsEnvs []models.DynamicSettingEnv
		sort.Strings(req.EnvironmentIds)
		for _, envID := range req.EnvironmentIds {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentSettingsVersion, envID)
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)

			dsEnvID, err := ids.IDFromRoot(existing.ProjectID, ids.DynamicSettingEnv)
			if err != nil {
				return err
			}
			proto, err := dynamic_settings.SerializeDefinition(ctx, req.Setting)
			if err != nil {
				return err
			}
			dsEnv := models.DynamicSettingEnv{
				Model:            models.Model{ID: dsEnvID},
				EnvironmentID:    ids.ID(envID),
				ProjectID:        existing.ProjectID,
				DynamicSettingID: existing.ID,
				// TODO: Need to ensure that we have a monotonic clock. We should take the max of
				// the environment's version and this. For now, let's use this.
				Version:     s.app.Clock().Now().UnixNano(),
				CreatedByID: user.ID,
				Proto:       proto,
			}
			dsEnvs = append(dsEnvs, dsEnv)
		}
		if err := tx.Save(&existing).Error; err != nil {
			return errors.WithStack(err)
		}

		if err := tx.Create(&dsEnvs).Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Errorf(codes.Internal, "could not update dynamic setting")
	}

	return &empty.Empty{}, nil
}

func (s *DashboardServer) DeleteDynamicSetting(ctx context.Context, req *pb_dashboard.DeleteDynamicSettingRequest) (*empty.Empty, error) {
	ds, err := s.authDynamicSetting(ctx, ids.ID(req.Id))
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
		proj, err := projects.GetProject(ctx, ds.ProjectID, tx, true)
		if err != nil {
			return err
		}
		sort.Slice(proj.Environments, func(i, j int) bool {
			return proj.Environments[i].ID < proj.Environments[j].ID
		})
		for _, env := range proj.Environments {
			// We must invalidate all the environment versions
			pending, err := s.app.KV().StartPending(ctx, kv.EnvironmentSettingsVersion, env.ID.String())
			if err != nil {
				return err
			}
			pendings = append(pendings, pending)
		}
		// Must pass since permissions are based on the project ID.
		if err := tx.Delete(&models.DynamicSetting{
			Model:     models.Model{ID: ids.ID(req.Id)},
			ProjectID: ds.ProjectID,
		}).Error; err != nil {
			return errors.WithStack(err)
		}
		// Delete it from all environments and all versions.
		if err := tx.Model(&models.DynamicSettingEnv{}).Where(
			"dynamic_setting_id = ? AND project_id = ?", req.Id, ds.ProjectID).Update("version", s.app.Clock().Now().UnixNano()).Error; err != nil {
			return errors.WithStack(err)
		}
		if err := tx.Delete(&models.DynamicSettingEnv{}, "dynamic_setting_id = ? AND project_id = ?", req.Id, ds.ProjectID).Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Errorf("%s\n", err)
		return nil, status.Error(codes.Internal, "could not delete dynamic setting")

	}
	return &empty.Empty{}, nil
}

func (s *DashboardServer) validateDynamicSetting(ds *pb_ds.DynamicSetting) error {
	if ds.Name == "" {
		return status.Error(codes.InvalidArgument, "Name is not specified")
	}
	if len(ds.Platforms) == 0 {
		return status.Error(codes.InvalidArgument, "No platform specified")
	}

	return nil
}
