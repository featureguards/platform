package environments

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/models"
	pb_project "platform/go/proto/project"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func Get(ctx context.Context, id ids.ID, db *gorm.DB) (*models.Environment, error) {
	var env models.Environment
	if err := db.WithContext(ctx).Where("id = ?", id).First(&env).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return &env, nil

}

func List(ctx context.Context, projectID ids.ID, db *gorm.DB) ([]models.Environment, error) {
	var envs []models.Environment
	if err := db.WithContext(ctx).Where("project_id = ?", projectID).Find(&envs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return envs, nil
}

func Pb(obj *models.Environment) (*pb_project.Environment, error) {
	res := &pb_project.Environment{
		Id:          string(obj.ID),
		CreatedAt:   timestamppb.New(obj.CreatedAt),
		Name:        obj.Name,
		Description: obj.Description,
		ProjectId:   string(obj.ProjectID),
	}
	return res, nil
}
