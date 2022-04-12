package projects

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	pb_project "stackv2/go/proto/project"
	"time"

	"github.com/golang/protobuf/ptypes"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	InviteExpiration = 24 * 7 * time.Hour
)

func GetProject(ctx context.Context, id ids.ID, db *gorm.DB) (*models.Project, error) {
	var project models.Project
	if err := db.WithContext(ctx).Where("id = ?", id).Preload("Environments").First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &project, nil

}

func PbMember(ctx context.Context, obj models.ProjectMember, ory *kratos.APIClient) (*pb_project.ProjectMember, error) {
	createdAt, err := ptypes.TimestampProto(obj.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: optimize by doing it concurrently. Ory doesn't support batching right now.
	identity, err := users.FetchIdentity(ctx, obj.User.OryID, ory)
	if err != nil {
		return nil, err
	}

	pbUser := users.Pb(identity, &obj.User)
	// Let's filter out some fields
	users.FilterPbUser(pbUser)

	res := &pb_project.ProjectMember{
		Id:        string(obj.ID),
		CreatedAt: createdAt,
		ProjectId: string(obj.ProjectID),
		User:      pbUser,
		Role:      obj.Role,
	}
	return res, nil
}

func PbProjectInvite(obj models.ProjectInvite) (*pb_project.ProjectInvite, error) {
	createdAt, err := ptypes.TimestampProto(obj.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	expiry := obj.CreatedAt.Add(InviteExpiration)
	expiresAt, err := ptypes.TimestampProto(expiry)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	status := obj.Status
	if status == pb_project.ProjectInvite_PENDING && time.Now().After(expiry) {
		status = pb_project.ProjectInvite_EXPIRED
	}

	invite := &pb_project.ProjectInvite{
		Id:          string(obj.ID),
		CreatedAt:   createdAt,
		ProjectId:   string(obj.ProjectID),
		ProjectName: obj.Project.Name,
		Email:       obj.Email,
		ExpiresAt:   expiresAt,
		Status:      status,
	}
	return invite, nil
}

func PbProject(proj models.Project) (*pb_project.Project, error) {
	createdAt, err := ptypes.TimestampProto(proj.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	envs := make([]*pb_project.Environment, len(proj.Environments))
	for i, env := range proj.Environments {
		envCreated, err := ptypes.TimestampProto(env.CreatedAt)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		envs[i] = &pb_project.Environment{
			Name:        env.Name,
			Description: env.Description,
			Id:          string(env.ID),
			CreatedAt:   envCreated,
		}
	}
	return &pb_project.Project{
		Name:         proj.Name,
		Id:           string(proj.ID),
		Description:  proj.Description,
		Environments: envs,
		CreatedAt:    createdAt,
	}, nil
}
