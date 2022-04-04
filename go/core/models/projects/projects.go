package projects

import (
	"context"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	pb_project "stackv2/go/proto/project"
	"time"

	"github.com/golang/protobuf/ptypes"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
)

const (
	InviteExpiration = 24 * 7 * time.Hour
)

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

	pbUser := users.PbUser(identity, &obj.User)
	// Let's filter out some fields
	users.FilterPbUser(pbUser)

	res := &pb_project.ProjectMember{
		Id:        string(obj.ID),
		CreatedAt: createdAt,
		ProjectId: obj.ProjectID,
		User:      pbUser,
		Role:      obj.Role,
	}
	return res, nil
}

func PbEnvironment(obj models.Environment) (*pb_project.Environment, error) {
	createdAt, err := ptypes.TimestampProto(obj.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := &pb_project.Environment{
		Id:          string(obj.ID),
		CreatedAt:   createdAt,
		Name:        obj.Name,
		Description: obj.Description,
		ProjectId:   obj.ProjectID,
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
		ProjectId:   obj.ProjectID,
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
