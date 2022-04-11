package dashboard

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/projects"
	"stackv2/go/core/models/users"
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

	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_ADMIN}); err != nil {
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

	if err := s.app.DB.WithContext(ctx).Create(&env).Error; err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create environment")
	}

	return projects.PbEnvironment(env)
}
func (s *DashboardServer) ListEnvironments(ctx context.Context, req *pb_dashboard.ListEnvironmentsRequest) (*pb_dashboard.ListEnvironmentsResponse, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}

	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	var envs []models.Environment
	if err := s.app.DB.WithContext(ctx).Where("project_id = ?", req.ProjectId).Find(&envs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no environments found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not list environments")
	}
	var pbEnvs []*pb_project.Environment
	for _, env := range envs {
		pbEnv, err := projects.PbEnvironment(env)
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
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	var env models.Environment
	if err := s.app.DB.WithContext(ctx).Where("id = ?", req.Id).First(&env).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no environment found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not find environment")

	}

	return projects.PbEnvironment(env)
}

func (s *DashboardServer) DeleteEnvironment(ctx context.Context, req *pb_dashboard.DeleteEnvironmentRequest) (*empty.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	if err := s.app.DB.WithContext(ctx).Delete(&models.Environment{Model: models.Model{ID: ids.ID(req.Id)}}).Error; err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not delete environment")
	}
	return &empty.Empty{}, nil

}
