package feature_toggles

import (
	"context"
	"fmt"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	"time"

	pb_ft "stackv2/go/proto/feature_toggle"
	pb_user "stackv2/go/proto/user"

	"github.com/golang/protobuf/ptypes"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const (
	InviteExpiration = 24 * 7 * time.Hour
)

type GetFTOpts struct {
	WithUser bool
}

func Get(ctx context.Context, id ids.ID, db *gorm.DB, opts GetFTOpts) (*models.FeatureToggle, error) {
	var ft models.FeatureToggle
	query := db.WithContext(ctx).Where("id = ?", id)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.First(&ft).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &ft, nil
}

func GetByName(ctx context.Context, projectID ids.ID, name string, db *gorm.DB, opts GetFTOpts) (*models.FeatureToggle, error) {
	var ft models.FeatureToggle
	query := db.WithContext(ctx).Where("project_id = ? AND name = ?", projectID, name)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.First(&ft).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return &ft, nil
}

func GetForEnv(ctx context.Context, envID, projectID ids.ID, db *gorm.DB, opts GetFTOpts) ([]models.FeatureToggle, error) {
	var fts []models.FeatureToggle
	query := db.WithContext(ctx).Where("environment_id = ? AND project_id = ?", envID, projectID)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.Find(&fts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err

	}
	return fts, nil
}

func GetLatestForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) (*models.FeatureToggleEnv, error) {
	var ftEnv models.FeatureToggleEnv
	if err := db.WithContext(ctx).Where("feature_toggle_id = ? AND environment_id = ?", id, envID).Order("version desc").Preload("CreatedBy").Preload("FeatureToggle").First(&ftEnv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err
	}

	return &ftEnv, nil
}

func ListLatestForEnv(ctx context.Context, envID ids.ID, db *gorm.DB) ([]models.FeatureToggleEnv, error) {
	type Pair struct {
		FeatureToggleID ids.ID
		Version         int64
	}
	var limitedFtEnvs []Pair
	if err := db.WithContext(ctx).Select("feature_toggle_id, MAX(version) as version").Where("environment_id = ?", envID).Group("feature_toggle_id").Table("feature_toggle_envs").Find(&limitedFtEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err
	}
	if len(limitedFtEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	var ftEnvs []models.FeatureToggleEnv
	where := make([][]interface{}, len(limitedFtEnvs))
	for i, pair := range limitedFtEnvs {
		where[i] = []interface{}{pair.FeatureToggleID, pair.Version}
	}
	if err := db.WithContext(ctx).Where("(feature_toggle_id, version) in ?", where).Order("created_at DESC").Preload("CreatedBy").Preload("FeatureToggle").Preload("FeatureToggle.CreatedBy").Find(&ftEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err
	}

	return ftEnvs, nil
}

func GetHistoryForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) ([]models.FeatureToggleEnv, error) {
	var ftEnvs []models.FeatureToggleEnv
	if err := db.WithContext(ctx).Where("environment_id = ? AND feature_toggle_id = ?", envID, id).Order("created_at DESC").Preload("CreatedBy").Preload("FeatureToggle").Find(&ftEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, err
	}
	if len(ftEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	return ftEnvs, nil
}

type PbOpts struct {
	FillUser bool
}

func Pb(ctx context.Context, ftEnv *models.FeatureToggleEnv, ory *kratos.APIClient, opts PbOpts) (*pb_ft.FeatureToggle, error) {
	ft := ftEnv.FeatureToggle
	createdAt, err := ptypes.TimestampProto(ft.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	updatedAt, err := ptypes.TimestampProto(ftEnv.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pbCreatedBy, pbUpdatedBy *pb_user.User
	if opts.FillUser {
		if ft.CreatedBy.OryID != "" {
			createdByIdenty, err := users.FetchIdentity(ctx, ft.CreatedBy.OryID, ory)
			if err != nil {
				return nil, err
			}
			pbCreatedBy = users.Pb(createdByIdenty, &ft.CreatedBy)
			// Let's filter out some fields
			users.LimitedPbUser(pbCreatedBy)
		}
		if ftEnv.CreatedBy.OryID != "" {
			updatedByIdenty, err := users.FetchIdentity(ctx, ftEnv.CreatedBy.OryID, ory)
			if err != nil {
				return nil, err
			}

			pbUpdatedBy = users.Pb(updatedByIdenty, &ftEnv.CreatedBy)
			users.LimitedPbUser(pbUpdatedBy)
		}
	}

	res := &pb_ft.FeatureToggle{
		Id:          string(ft.ID),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Name:        ft.Name,
		Description: ft.Description,
		ToggleType:  ft.Type,
		Version:     ftEnv.Version,
		Enabled:     ftEnv.Enabled,
		Platforms:   []pb_ft.Platform_Type{pb_ft.Platform_DEFAULT},
		ProjectId:   string(ft.ProjectID),
		CreatedBy:   pbCreatedBy,
		UpdatedBy:   pbUpdatedBy,
	}
	if ft.IsMobile {
		res.Platforms = append(res.Platforms, pb_ft.Platform_MOBILE)
	}
	if ft.IsWeb {
		res.Platforms = append(res.Platforms, pb_ft.Platform_WEB)
	}
	pb, err := ftEnv.ProtoMessage(ft.Type)
	if err != nil {
		return nil, err
	}
	if err := fillProto(res, pb); err != nil {
		return nil, err
	}
	return res, nil
}

func fillProto(ft *pb_ft.FeatureToggle, pb proto.Message) error {
	switch ft.ToggleType {
	case pb_ft.FeatureToggle_EXPERIMENT:
		ft.FeatureDefinition = &pb_ft.FeatureToggle_Experiment{Experiment: pb.(*pb_ft.ExperimentFeature)}
	case pb_ft.FeatureToggle_ON_OFF:
		ft.FeatureDefinition = &pb_ft.FeatureToggle_OnOff{OnOff: pb.(*pb_ft.OnOffFeature)}
	case pb_ft.FeatureToggle_PERCENTAGE:
		ft.FeatureDefinition = &pb_ft.FeatureToggle_Percentage{Percentage: pb.(*pb_ft.PercentageFeature)}
	case pb_ft.FeatureToggle_PERMISSION:
		ft.FeatureDefinition = &pb_ft.FeatureToggle_Permission{Permission: pb.(*pb_ft.PermissionFeature)}
	default:
		return errors.WithStack(fmt.Errorf("unknown toggle type %v", ft.ToggleType))
	}
	return nil
}
