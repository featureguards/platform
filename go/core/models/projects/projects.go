package projects

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	pb_project "stackv2/go/proto/project"
	"time"

	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	InviteExpiration = 24 * 7 * time.Hour
)

func GetProject(ctx context.Context, id ids.ID, db *gorm.DB, lock bool) (*models.Project, error) {
	var project models.Project
	q := db.WithContext(ctx)
	if lock {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.Where("id = ?", id).Preload("Environments").First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &project, nil

}

func PbMember(ctx context.Context, obj models.ProjectMember, ory *kratos.APIClient) (*pb_project.ProjectMember, error) {
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
		CreatedAt: timestamppb.New(obj.CreatedAt),
		ProjectId: string(obj.ProjectID),
		User:      pbUser,
		Role:      obj.Role,
	}
	return res, nil
}

func PbProjectInvite(obj models.ProjectInvite) (*pb_project.ProjectInvite, error) {
	expiry := obj.CreatedAt.Add(InviteExpiration)
	status := obj.Status
	if status == pb_project.ProjectInvite_PENDING && time.Now().After(expiry) {
		status = pb_project.ProjectInvite_EXPIRED
	}

	invite := &pb_project.ProjectInvite{
		Id:          string(obj.ID),
		CreatedAt:   timestamppb.New(obj.CreatedAt),
		ProjectId:   string(obj.ProjectID),
		ProjectName: obj.Project.Name,
		Email:       obj.Email,
		ExpiresAt:   timestamppb.New(expiry),
		Status:      status,
	}
	return invite, nil
}

func PbProject(proj models.Project) (*pb_project.Project, error) {
	envs := make([]*pb_project.Environment, len(proj.Environments))
	for i, env := range proj.Environments {
		envs[i] = &pb_project.Environment{
			Name:        env.Name,
			Description: env.Description,
			Id:          string(env.ID),
			CreatedAt:   timestamppb.New(env.CreatedAt),
		}
	}
	return &pb_project.Project{
		Name:         proj.Name,
		Id:           string(proj.ID),
		Description:  proj.Description,
		Environments: envs,
		CreatedAt:    timestamppb.New(proj.CreatedAt),
	}, nil
}
