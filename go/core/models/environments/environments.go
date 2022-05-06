package environments

import (
	"context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	pb_project "stackv2/go/proto/project"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func Get(ctx context.Context, id ids.ID, db *gorm.DB) (*models.Environment, error) {
	var env models.Environment
	if err := db.WithContext(ctx).Where("id = ?", id).First(&env).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &env, nil

}

func List(ctx context.Context, projectID ids.ID, db *gorm.DB) ([]models.Environment, error) {
	var envs []models.Environment
	if err := db.WithContext(ctx).Where("project_id = ?", projectID).Find(&envs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

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
