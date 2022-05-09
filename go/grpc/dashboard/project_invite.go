package dashboard

import (
	"context"
	"net/http"
	"platform/go/core/app_context"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/models/projects"
	"platform/go/core/models/users"
	"platform/go/core/ory"
	pb_dashboard "platform/go/proto/dashboard"
	"platform/go/proto/project"
	pb_project "platform/go/proto/project"
	"strings"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type listInvitesReq struct {
	userID    string
	projectID string
}

func (s *DashboardServer) CreateProjectInvite(ctx context.Context, req *pb_dashboard.CreateProjectInviteRequest) (*empty.Empty, error) {
	// We validate here
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is not specified")
	}
	if req.Invite.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is not specified")
	}
	if req.Invite.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first name is not specified")
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
	email := strings.ToLower(strings.TrimSpace(req.Invite.Email))
	if err := s.DB(ctx).Where("project_id = ? AND email = ? AND created_at > ?", req.ProjectId, email, time.Now().Add(-projects.InviteExpiration)).Find(&invites).Error; err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}
	if len(invites) > 0 {
		// We will extend their expiry
		for i := range invites {
			invites[i].ExpiresAt = time.Now().Add(projects.InviteExpiration)
		}
		if err := s.DB(ctx).Save(invites).Error; err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not create project invite")
		}
		return &empty.Empty{}, nil
	}
	// No such invite. Let's create one.
	id, err := ids.IDFromRoot(ids.ID(req.ProjectId), ids.ProjectInvite)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}
	// Let's find out if need to create a new identity first.
	identity, res, err := s.app.Ory.AdminApi().V0alpha2Api.AdminCreateIdentity(ctx).AdminCreateIdentityBody(*kratos.NewAdminCreateIdentityBody(ory.Person, map[string]interface{}{"email": email, "first_name": strings.TrimSpace(req.Invite.FirstName)})).Execute()
	if err != nil {
		if res.StatusCode == http.StatusConflict {
			// TODO: Support this. We need to list all idenitiess. Skip for now.
		}
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}

	invite := models.ProjectInvite{
		Model:     models.Model{ID: id},
		OryID:     identity.Id,
		Status:    pb_project.ProjectInvite_PENDING,
		ProjectID: ids.ID(req.ProjectId),
		ExpiresAt: time.Now().Add(projects.InviteExpiration),
	}
	if res := s.DB(ctx).Save(&invite); res.Error != nil {
		log.Error(errors.WithStack(res.Error))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}
	body := kratos.NewAdminCreateSelfServiceRecoveryLinkBody("8ea268a2-abc1-433e-b5b8-c28d6545695d")
	link, _, err := s.app.Ory.AdminApi().V0alpha2Api.AdminCreateSelfServiceRecoveryLink(context.Background()).AdminCreateSelfServiceRecoveryLinkBody(*body).Execute()
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}
	log.Infoln(link.RecoveryLink)
	return &empty.Empty{}, nil
}
func (s *DashboardServer) ListProjectInvites(ctx context.Context, req *pb_dashboard.ListProjectInvitesRequest) (*pb_project.ProjectInvites, error) {
	return s.listProjectOrUserInvites(ctx, listInvitesReq{projectID: req.ProjectId})
}

func (s *DashboardServer) ListUserInvites(ctx context.Context, req *pb_dashboard.ListUserInvitesRequest) (*pb_project.ProjectInvites, error) {
	return s.listProjectOrUserInvites(ctx, listInvitesReq{userID: req.UserId})
}

func (s *DashboardServer) listProjectOrUserInvites(ctx context.Context, req listInvitesReq) (*pb_project.ProjectInvites, error) {
	// We list either by project ID or userID
	var invites []models.ProjectInvite
	if req.projectID != "" && req.userID != "" {
		return nil, status.Error(codes.InvalidArgument, "either project_id or user_id must be set")
	}
	if req.projectID == "" && req.userID == "" {
		return nil, status.Error(codes.InvalidArgument, "one of project_id or user_id must be set")
	}
	if req.projectID != "" {
		if _, err := s.authProject(ctx, ids.ID(req.projectID), adminOnly); err != nil {
			return nil, err
		}
		// user is a member of the project. Return all invites
		if err := s.DB(ctx).Where("project_id = ?", req.projectID).Find(&invites).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// We're good
			} else {
				log.Error(errors.WithStack(err))
				return nil, status.Error(codes.Internal, "could not list project invites")
			}

		}
	} else if req.userID != "" {
		var userID ids.ID
		user, err := users.FetchUserForSession(ctx, s.app.DB)
		if err != nil {
			return nil, status.Error(codes.NotFound, "no user for session")
		}
		if req.userID == users.Me {
			userID = user.ID
		}
		// Access check
		if userID != user.ID {
			return nil, status.Error(codes.NotFound, "no project invite exists for user")
		}
		if err := s.DB(ctx).Where("ory_id = ?", user.OryID).Find(&invites).Error; err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not list project invites")
		}
	}
	var pb_invites []*pb_project.ProjectInvite
	for _, invite := range invites {
		identity, err := users.FetchIdentity(ctx, invite.OryID, s.app.Ory.Api())
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not list project invites")
		}

		pb_invite, err := projects.PbProjectInvite(invite, identity)
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
	res := s.DB(ctx).Where("id = ?", req.Id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "no project invite exists")
	}

	// Authorize based on OryId only
	if invite.OryID != user.OryID {
		return nil, status.Error(codes.NotFound, "no project invite exists")
	}

	return projects.PbProjectInvite(invite, &session.Identity)
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
	if res := s.DB(ctx).Model(&models.ProjectInvite{}).Where("id = ?", req.Invite.Id).Select("status").Updates(fields); res.Error != nil {
		log.Error(errors.WithStack(res.Error))
		return nil, status.Error(codes.Internal, "could not update project invite")
	}
	return s.GetProjectInvite(ctx, &pb_dashboard.GetProjectInviteRequest{Id: req.Invite.Id})
}
