package projects

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/models/users"
	"platform/go/core/ory"
	pb_project "platform/go/proto/project"
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

func GetProjectMember(ctx context.Context, id ids.ID, db *gorm.DB) (*models.ProjectMember, error) {
	var projectMember models.ProjectMember
	q := db.WithContext(ctx)

	if err := q.Where("id = ?", id).First(&projectMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &projectMember, nil
}

func GetProjectInvite(ctx context.Context, id ids.ID, db *gorm.DB) (*models.ProjectInvite, error) {
	var invite models.ProjectInvite

	if err := db.WithContext(ctx).Where("id = ?", id).Preload("Project").Find(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err
	}

	return &invite, nil
}

func PbMember(ctx context.Context, obj models.ProjectMember, ory *ory.Ory) (*pb_project.ProjectMember, error) {
	// TODO: optimize by doing it concurrently. Ory doesn't support batching right now.
	identity, err := users.FetchIdentity(ctx, obj.User.OryID, ory.Api())
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

func PbProjectInvite(obj models.ProjectInvite, identity *kratos.Identity) (*pb_project.ProjectInvite, error) {
	traits := ory.Traits(identity.Traits.(map[string]interface{}))
	invite := &pb_project.ProjectInvite{
		Id:          string(obj.ID),
		CreatedAt:   timestamppb.New(obj.CreatedAt),
		ProjectId:   string(obj.ProjectID),
		ProjectName: obj.Project.Name,
		Status:      obj.DerivedStatus(),
		Email:       identity.VerifiableAddresses[0].Value,
		ExpiresAt:   timestamppb.New(obj.ExpiresAt),
		FirstName:   traits.FirstName(),
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
