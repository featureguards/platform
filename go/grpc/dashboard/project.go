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

func (s *DashboardServer) CreateProject(ctx context.Context, req *pb_dashboard.CreateProjectRequest) (*pb_project.Project, error) {
	// We validate here
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "project name is not specified")
	}

	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	projectID, err := ids.RandomID(ids.Project)
	if err != nil {
		err := errors.WithStack(err)
		log.Error(err)
		return nil, status.Error(codes.Internal, "could not create project")
	}
	envs := make([]models.Environment, len(req.Environments))
	for i, env := range req.Environments {
		envs[i] = models.Environment{
			Name:        env.Name,
			Description: env.Description,
		}
		envID, err := ids.IDFromRoot(projectID, ids.Environment)
		if err != nil {
			err := errors.WithStack(err)
			log.Error(err)
			return nil, status.Errorf(codes.Internal, "could not retrieve user")
		}
		envs[i].ID = envID

	}
	// Add user as member by default.
	projectMemberID, err := ids.IDFromRoot(projectID, ids.ProjectMember)
	if err != nil {
		err := errors.WithStack(err)
		log.Error(err)
		return nil, status.Errorf(codes.Internal, "could not retrieve user")
	}
	projectMember := models.ProjectMember{
		Model:  models.Model{ID: projectMemberID},
		UserID: string(user.ID), Role: pb_project.Project_ADMIN,
	}
	project := models.Project{
		Model:          models.Model{ID: projectID},
		Name:           req.Name,
		Description:    req.Description,
		Environments:   envs,
		OwnerID:        string(user.ID),
		ProjectMembers: []models.ProjectMember{projectMember},
	}
	if err := s.app.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&project)
		if res.Error != nil {
			log.Error(errors.WithStack(res.Error))
			return status.Error(codes.Internal, "error creating project")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	pb, err := projects.PbProject(project)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error creating project")
	}
	return pb, nil
}

func (s *DashboardServer) ListProjects(ctx context.Context, req *pb_dashboard.ListProjectsRequest) (*pb_dashboard.ListProjectsResponse, error) {
	var members []models.ProjectMember
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.app.DB.WithContext(ctx).Where("user_id = ?", user.ID).Preload("Project").Find(&members).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no projects found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not retrieve projects")
	}

	pbProjects := make([]*pb_project.Project, 0, len(members))
	// Fetch the projects
	for _, member := range members {
		pb, err := projects.PbProject(member.Project)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not retrieve projects")
		}
		pbProjects = append(pbProjects, pb)
	}
	return &pb_dashboard.ListProjectsResponse{
		Projects: pbProjects,
	}, nil
}

func (s *DashboardServer) GetProject(ctx context.Context, req *pb_dashboard.GetProjectRequest) (*pb_project.Project, error) {
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	return s.getProjectForUser(ctx, user.ID, ids.ID(req.Id))
}

func (s *DashboardServer) DeleteProject(ctx context.Context, req *pb_dashboard.DeleteProjectRequest) (*empty.Empty, error) {
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.Id), []pb_project.Project_Role{pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	if err := s.app.DB.WithContext(ctx).Delete(&models.Project{
		Model: models.Model{ID: ids.ID(req.Id)},
	}).Error; err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not delete project")
	}
	return &empty.Empty{}, nil
}

func (s *DashboardServer) ListProjectMembers(ctx context.Context, req *pb_dashboard.ListProjectMembersRequest) (*pb_project.ProjectMembers, error) {
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	var members []models.ProjectMember
	if err := s.app.DB.WithContext(ctx).Where("project_id = ?", req.ProjectId).Find(&members).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(members) <= 0 {
			// Unauthorized or project not found
			return nil, status.Error(codes.NotFound, "no project found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not retrieve project members")
	}

	var pbMembers []*pb_project.ProjectMember
	for _, member := range members {
		pbMember, err := projects.PbMember(ctx, member, s.app.Ory)
		if err != nil {
			log.Error(err)
			return nil, status.Error(codes.Internal, "could not retrieve project member")
		}
		pbMembers = append(pbMembers, pbMember)
	}

	return &pb_project.ProjectMembers{
		Members: pbMembers,
	}, nil
}

func (s *DashboardServer) validateMembership(ctx context.Context, userID, projectID ids.ID, roles []pb_project.Project_Role) error {
	rolesMap := make(map[pb_project.Project_Role]struct{}, len(roles))
	for _, role := range roles {
		rolesMap[role] = struct{}{}
	}
	var members []models.ProjectMember
	if err := s.app.DB.WithContext(ctx).Where("user_id = ? AND project_id = ?", userID, projectID).Find(&members).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(members) <= 0 {
			// Unauthorized or project not found
			return status.Error(codes.NotFound, "no project found")
		}
		log.Error(errors.WithStack(err))
		return status.Error(codes.Internal, "could not retrieve project")
	}
	if len(members) <= 0 {
		// Unauthorized or project not found
		return status.Error(codes.NotFound, "no project found")
	}
	for _, member := range members {
		if _, ok := rolesMap[member.Role]; ok {
			return nil
		}
	}

	// Don't leak permissions. Hence, return a 404 even though this is a 403
	return status.Error(codes.NotFound, "no project found")
}

func (s *DashboardServer) getProjectForUser(ctx context.Context, userID, id ids.ID) (*pb_project.Project, error) {
	var members []models.ProjectMember
	if err := s.app.DB.WithContext(ctx).Where("user_id = ? AND project_id = ?", userID, id).Find(&members).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Unauthorized or project not found
			return nil, status.Error(codes.NotFound, "no project found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not retrieve project")
	}

	if len(members) <= 0 {
		// Unauthorized or project not found
		return nil, status.Error(codes.NotFound, "no project found")
	}

	var project models.Project
	if err := s.app.DB.WithContext(ctx).Where("id = ?", id).Preload("Environments").First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "no projects found")
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not retrieve projects")

	}
	return projects.PbProject(project)
}
