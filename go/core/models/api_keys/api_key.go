package api_keys

import (
	"context"
	"platform/go/core/ids"
	"platform/go/core/models"
	pb_project "platform/go/proto/project"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func Get(ctx context.Context, id ids.ID, db *gorm.DB) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	if err := db.WithContext(ctx).Where("id = ?", id).Find(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return &apiKey, nil
}

func List(ctx context.Context, environmentID ids.ID, db *gorm.DB) ([]models.ApiKey, error) {
	var apiKeys []models.ApiKey
	if err := db.WithContext(ctx).Where("environment_id = ?", environmentID).Find(&apiKeys).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return apiKeys, nil
}

func Pb(obj *models.ApiKey) (*pb_project.ApiKey, error) {
	res := &pb_project.ApiKey{
		Id:            string(obj.ID),
		CreatedAt:     timestamppb.New(obj.CreatedAt),
		Name:          obj.Name,
		ProjectId:     string(obj.ProjectID),
		EnvironmentId: string(obj.EnvironmentID),
		Key:           string(obj.Key),
	}
	if !obj.ExpiresAt.IsZero() {
		res.ExpiresAt = timestamppb.New(obj.ExpiresAt)
	}
	return res, nil
}
