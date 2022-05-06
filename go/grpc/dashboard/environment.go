package dashboard

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/environments"
	"stackv2/go/core/models/feature_toggles"
	"stackv2/go/core/models/projects"
	pb_dashboard "stackv2/go/proto/dashboard"
	pb_project "stackv2/go/proto/project"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DashboardServer) CreateEnvironment(ctx context.Context, req *pb_dashboard.CreateEnvironmentRequest) (*pb_project.Environment, error) {
	// We validate here
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}

	if _, err := s.authProject(ctx, ids.ID(req.ProjectId), adminOnly); err != nil {
		return nil, err
	}

	id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.Environment)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "could not create environment")
	}

	env := models.Environment{
		Model:       models.Model{ID: id},
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   ids.ID(req.ProjectId),
	}

	// Take a lock
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := projects.GetProject(ctx, ids.ID(req.ProjectId), tx, true); err != nil {
			return err
		}
		if err := tx.Create(&env).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Error(codes.Internal, "could not create environment")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return environments.Pb(&env)
}
func (s *DashboardServer) ListEnvironments(ctx context.Context, req *pb_dashboard.ListEnvironmentsRequest) (*pb_dashboard.ListEnvironmentsResponse, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}

	if _, err := s.authProject(ctx, ids.ID(req.ProjectId), allRoles); err != nil {
		return nil, err
	}

	var envs []models.Environment
	if err := s.DB(ctx).Where("project_id = ?", req.ProjectId).Find(&envs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no environments found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not list environments")
	}
	var pbEnvs []*pb_project.Environment
	for _, env := range envs {
		pbEnv, err := environments.Pb(&env)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not list environments")
		}
		pbEnvs = append(pbEnvs, pbEnv)
	}
	return &pb_dashboard.ListEnvironmentsResponse{
		Environments: pbEnvs,
	}, nil
}

func (s *DashboardServer) GetEnvironment(ctx context.Context, req *pb_dashboard.GetEnvironmentRequest) (*pb_project.Environment, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}
	env, err := s.authEnvironment(ctx, ids.ID(req.Id))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not find environment")
	}

	return environments.Pb(env)
}

func (s *DashboardServer) CloneEnvironment(ctx context.Context, req *pb_dashboard.CloneEnvironmentRequest) (*pb_project.Environment, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	existing, err := s.authEnvironment(ctx, ids.ID(req.Id))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not find environment")
	}

	// Clone
	var env models.Environment
	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the project
		if _, err := projects.GetProject(ctx, existing.ProjectID, tx, true); err != nil {
			return err
		}
		id, err := ids.IDFromRoot(ids.ID(existing.ProjectID), ids.Environment)
		if err != nil {
			log.Error(err)
			return status.Error(codes.Internal, "could not create environment")
		}

		// Create a new environment
		env = models.Environment{
			Model:       models.Model{ID: id},
			Name:        req.Name,
			Description: req.Description,
			ProjectID:   ids.ID(existing.ProjectID),
		}
		if err := tx.Create(&env).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Error(codes.Internal, "could not clone environment")
		}

		ftEnvs, err := feature_toggles.ListLatestForEnv(ctx, ids.ID(env.ID), tx)
		if err != nil {
			if err == models.ErrNotFound {
				return status.Errorf(codes.NotFound, "feature toggles not found")
			}
			return status.Errorf(codes.Internal, "could not clone environment")
		}

		for _, ftEnv := range ftEnvs {
			ftEnv.EnvironmentID = env.ID
			ftEnv.Environment = env
		}
		if err := tx.Save(&ftEnvs).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Errorf(codes.Internal, "could not clone environment")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return environments.Pb(&env)
}

func (s *DashboardServer) DeleteEnvironment(ctx context.Context, req *pb_dashboard.DeleteEnvironmentRequest) (*empty.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}
	env, err := s.authEnvironment(ctx, ids.ID(req.Id))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		return nil, status.Error(codes.Internal, "could not find environment")
	}

	if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the project
		_, err := projects.GetProject(ctx, env.ProjectID, tx, true)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return status.Error(codes.NotFound, "no projec found")
			}
			return status.Error(codes.Internal, "could not delete environment")
		}

		// Can't delete the last environment
		envs, err := environments.List(ctx, env.ProjectID, tx)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return status.Error(codes.NotFound, "no environments found")
			}
			return status.Error(codes.Internal, "could not delete environment")
		}

		if len(envs) <= 1 {
			return status.Error(codes.InvalidArgument, "must have at least one environment")
		}

		if err := tx.Delete(&models.Environment{Model: models.Model{ID: ids.ID(req.Id)}}).Error; err != nil {
			log.Error(errors.WithStack(err))
			return status.Error(codes.Internal, "could not delete environment")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil

}
