package dashboard

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/feature_toggles"
	"stackv2/go/core/models/users"
	pb_project "stackv2/go/proto/project"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	allRoles  = []pb_project.Project_Role{pb_project.Project_MEMBER, pb_project.Project_ADMIN}
	adminOnly = []pb_project.Project_Role{pb_project.Project_ADMIN}
)

func (s *DashboardServer) authFeatureToggle(ctx context.Context, id ids.ID) (*models.FeatureToggle, error) {
	ft, err := feature_toggles.Get(ctx, id, s.app.DB, feature_toggles.GetFTOpts{})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no feature toggle found")
		}
		return nil, status.Errorf(codes.Internal, "could not get feature toggle")
	}
	if _, err := s.authProject(ctx, ids.ID(ft.ProjectID), allRoles); err != nil {
		return nil, err
	}

	return ft, nil
}

func (s *DashboardServer) authProject(ctx context.Context, id ids.ID, roles []pb_project.Project_Role) (*models.User, error) {
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	if err := s.validateMembership(ctx, user.ID, id, roles); err != nil {
		return nil, err
	}

	return user, nil
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
