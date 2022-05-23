package dashboard

import (
	"context"
	"net/http"
	"platform/go/core/app_context"
	"platform/go/core/ids"
	"platform/go/core/mail/templates"
	"platform/go/core/models"
	"platform/go/core/models/projects"
	"platform/go/core/models/users"
	"platform/go/core/ory"
	pb_dashboard "platform/go/proto/dashboard"
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

	email := strings.ToLower(strings.TrimSpace(req.Invite.Email))
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
			// TODO: Optimize by changing Kratos API to accept an email
			identities, _, err := s.app.Ory.AdminApi().V0alpha2Api.AdminListIdentities(ctx).Execute()
			if err != nil {
				log.Error(errors.WithStack(err))
				return nil, status.Error(codes.Internal, "could not create project invite")
			}
		loop:
			for i, ident := range identities {
				for _, addr := range ident.VerifiableAddresses {
					if addr.Value == email {
						identity = &identities[i]
						// Fetch the invite
						var invite models.ProjectInvite
						if err := s.DB(ctx).First(&invite).Where("ory_id = ?", identity.Id).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								log.Error(errors.WithStack(err))
								return nil, status.Error(codes.Internal, "could not create project invite")
							}
						} else {
							id = invite.ID
						}
						break loop
					}
				}
			}
		} else {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not create project invite")
		}
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
	body := kratos.NewAdminCreateSelfServiceRecoveryLinkBody(identity.Id)
	link, _, err := s.app.Ory.AdminApi().V0alpha2Api.AdminCreateSelfServiceRecoveryLink(context.Background()).AdminCreateSelfServiceRecoveryLinkBody(*body).Execute()
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}

	proj, err := s.getProjectForUser(ctx, user.ID, ids.ID(req.ProjectId))
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not create project invite")
	}
	session, _ := app_context.SessionFromContext(ctx)
	traits := ory.Traits(session.Identity.Traits.(map[string]interface{}))

	if err := s.app.Mail.Send(ctx, templates.NewProjectInvitationTemplate(&templates.ProjectInvite{
		FirstName: req.Invite.FirstName,
		Email:     req.Invite.Email,
		Link:      link.RecoveryLink,
		Project:   proj.Name,
		Sender:    traits.FirstName(),
	})); err != nil {
		log.Errorf("%+v", errors.WithStack(err))
		return nil, status.Error(codes.Internal, "could not send email for project invite")
	}

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
		if err := s.DB(ctx).Where("project_id = ?", req.projectID).Preload("Project").Find(&invites).Error; err != nil {
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
		if err := s.DB(ctx).Where("ory_id = ?", user.OryID).Preload("Project").Find(&invites).Error; err != nil {
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

func (s *DashboardServer) GetProjectInvite(ctx context.Context, req *pb_dashboard.GetProjectInviteRequest) (*pb_project.ProjectInvite, error) {
	session, ok := app_context.SessionFromContext(ctx)
	if !ok {
		return nil, models.ErrNoSession
	}
	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	invite, err := projects.GetProjectInvite(ctx, ids.ID(req.Id), s.DB(ctx))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no project invite exists")
		}
		return nil, status.Error(codes.Internal, "could not get project invite")
	}

	// Authorize based on OryId only
	if invite.OryID != user.OryID {
		return nil, status.Error(codes.NotFound, "no project invite exists")
	}

	return projects.PbProjectInvite(*invite, &session.Identity)
}

func (s *DashboardServer) UpdateProjectInvite(ctx context.Context, req *pb_dashboard.UpdateProjectInviteRequest) (*pb_project.ProjectInvite, error) {
	if req.Invite == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "no invite specified")
	}

	user, err := users.FetchUserForSession(ctx, s.app.DB)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no user for session")
	}

	// This checks that the user has access.
	existing, err := projects.GetProjectInvite(ctx, ids.ID(req.Id), s.DB(ctx))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no project invite exists")
		}
		return nil, status.Error(codes.Internal, "could not get project invite")
	}

	// Authorize based on OryId only
	if existing.OryID != user.OryID {
		return nil, status.Error(codes.NotFound, "no project invite exists")
	}

	if existing.Status == pb_project.ProjectInvite_EXPIRED || existing.Status == pb_project.ProjectInvite_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "status must be valid")
	}

	// For now, can only change status
	if req.Invite.Status == pb_project.ProjectInvite_ACCEPTED {
		fields := models.FieldsFromPb(req.Invite.ProtoReflect())
		if err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.ProjectInvite{}).Where("id = ?", req.Id).Select("status").Updates(fields).Error; err != nil {
				return errors.WithStack(err)
			}

			var projectMember models.ProjectMember
			err = tx.Where("user_id = ?", user.ID).First(&projectMember).Error
			if err == nil {
				return nil
			}
			if err != nil && err != gorm.ErrRecordNotFound {
				return errors.WithStack(err)
			}
			// Add user to project members
			projectMemberID, err := ids.IDFromRoot(ids.ID(existing.ProjectID), ids.ProjectMember)
			if err != nil {
				return err
			}
			projectMember = models.ProjectMember{
				Model:  models.Model{ID: projectMemberID},
				UserID: user.ID, Role: pb_project.Project_MEMBER,
				ProjectID: existing.ProjectID,
			}
			if err := tx.Save(&projectMember).Error; err != nil {
				return errors.WithStack(err)
			}

			return nil
		}); err != nil {
			log.Error(errors.WithStack(err))
			return nil, status.Error(codes.Internal, "could not update project invite")
		}
	}
	return s.GetProjectInvite(ctx, &pb_dashboard.GetProjectInviteRequest{Id: req.Id})
}
