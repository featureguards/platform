package dashboard

import (
	"context"
	"stackv2/go/core/app_context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/projects"
	"stackv2/go/core/models/users"
	pb_dashboard "stackv2/go/proto/dashboard"
	"stackv2/go/proto/project"
	pb_project "stackv2/go/proto/project"
	"strings"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *DashboardServer) CreateProjectInvite(ctx context.Context, req *pb_dashboard.ProjectInviteRequest) (*empty.Empty, error) {
	// We validate here
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is not specified")
	}

	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_ADMIN}); err != nil {
		return nil, err
	}

	// See if invite already exists and hasn't expired.
	var invites []models.ProjectInvite
	res := s.app.DB.WithContext(ctx).Where("project_id = ? AND email = ? AND created_at > ?", req.ProjectId, strings.ToLower(req.Email), time.Now().Add(-projects.InviteExpiration)).Find(&invites)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// No such invite. Let's create one.
		id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.ProjectInvite)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not create project invite")
		}
		invite := models.ProjectInvite{
			Model:  models.Model{ID: id},
			Email:  strings.ToLower(req.Email),
			Status: pb_project.ProjectInvite_PENDING,
		}
		if res := s.app.DB.Create(&invite); res.Error != nil {
			log.Error(errors.WithStack(res.Error))
			return nil, status.Error(codes.Internal, "could not create project invite")
		}
		return nil, status.Error(codes.NotFound, "no projects found")
	}
	return nil, status.Errorf(codes.AlreadyExists, "a project invite already exists")
}
func (s *DashboardServer) ListProjectInvites(ctx context.Context, req *pb_dashboard.ListProjectInvitesRequest) (*pb_project.ProjectInvites, error) {
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}
	// We list either by project ID or userID
	var invites []models.ProjectInvite
	if req.ProjectId != "" && req.UserId != "" {
		return nil, status.Error(codes.InvalidArgument, "either project_id or user_id must be set")
	}
	if req.ProjectId == "" && req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "one of project_id or user_id must be set")
	}
	if req.ProjectId != "" {
		if err := s.validateMembership(ctx, user.ID, ids.ID(req.ProjectId), []pb_project.Project_Role{pb_project.Project_ADMIN}); err != nil {
			return nil, err
		}
		// user is a member of the project. Return all invites
		res := s.app.DB.WithContext(ctx).Where("project_id = ?", req.ProjectId).Find(&invites)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// We're good
		} else if res.Error != nil {
			log.Error(errors.WithStack(res.Error))
			return nil, status.Error(codes.Internal, "could not list project invites")
		}
	} else if req.UserId != "" {
		res := s.app.DB.WithContext(ctx).Where("user_id = ?", req.UserId).Find(&invites)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// We're good
		} else if res.Error != nil {
			log.Error(errors.WithStack(res.Error))
			return nil, status.Error(codes.Internal, "could not list project invites")
		}
	}
	var pb_invites []*pb_project.ProjectInvite
	for _, invite := range invites {
		pb_invite, err := projects.PbProjectInvite(invite)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not list project invites")
		}
		pb_invites = append(pb_invites, pb_invite)
	}
	return &pb_project.ProjectInvites{
		Invites: pb_invites,
	}, nil
}

func (s *DashboardServer) GetProjectInvite(ctx context.Context, req *pb_dashboard.GetProjectInviteRequest) (*project.ProjectInvite, error) {
	session, ok := app_context.SessionFromContext(ctx)
	if !ok {
		return nil, models.ErrNoSession
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	var invite models.ProjectInvite
	res := s.app.DB.WithContext(ctx).Where("id = ?", req.Id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no project invite exists")
	}

	// Authorize based on validated user emails only
	pbUser := users.PbUser(&session.Identity, user)
	for _, addr := range pbUser.Addresses {
		if addr.Verified && strings.ToLower(addr.Address) == strings.ToLower(invite.Email) {
			return projects.PbProjectInvite(invite)
		}
	}
	return nil, status.Error(codes.NotFound, "no project invite exists")
}
func (s *DashboardServer) UpdateProjectInvite(ctx context.Context, req *pb_dashboard.UpdateProjectInviteRequest) (*project.ProjectInvite, error) {
	if req.Invite == nil || req.Invite.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "no invite specified")
	}

	if req.Invite.Status == pb_project.ProjectInvite_EXPIRED || req.Invite.Status == pb_project.ProjectInvite_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "status must be valid")
	}

	// This checks that the user has access.
	if _, err := s.GetProjectInvite(ctx, &pb_dashboard.GetProjectInviteRequest{Id: req.Invite.Id}); err != nil {
		return nil, err
	}

	fields := models.FieldsFromPb(req.Invite.ProtoReflect())
	if res := s.app.DB.WithContext(ctx).Model(&models.ProjectInvite{}).Where("id = ?", req.Invite.Id).Select("status").Updates(fields); res.Error != nil {
		log.Error(errors.WithStack(res.Error))
		return nil, status.Error(codes.Internal, "could not update project invite")
	}
	return s.GetProjectInvite(ctx, &pb_dashboard.GetProjectInviteRequest{Id: req.Invite.Id})
}
