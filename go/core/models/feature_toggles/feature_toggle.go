package feature_toggles

import (
	"context"
	"fmt"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/models/users"
	"platform/go/core/ory"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
	pb_user "github.com/featureguards/featureguards-go/v2/proto/user"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
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
		return nil, errors.WithStack(err)

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
		return nil, errors.WithStack(err)

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
		return nil, errors.WithStack(err)

	}
	return fts, nil
}

func GetLatestForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) (*models.FeatureToggleEnv, error) {
	var ftEnv models.FeatureToggleEnv
	if err := db.WithContext(ctx).Where("feature_toggle_id = ? AND environment_id = ?", id, envID).Order("version desc").Preload("CreatedBy").Preload("FeatureToggle").First(&ftEnv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &ftEnv, nil
}

func MaxVersionForEnv(ctx context.Context, envID ids.ID, db *gorm.DB) (int64, error) {
	var version int64
	if err := db.WithContext(ctx).Select("MAX(version) as version").Where("environment_id = ?", envID).Table("feature_toggle_envs").Find(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, models.ErrNotFound
		}
		return -1, errors.WithStack(err)
	}
	return version, nil
}

type listOptions struct {
	start   int64
	end     int64
	deleted bool
}

type ListOptions func(o *listOptions) error

func WithStartVersion(v int64) ListOptions {
	return func(o *listOptions) error {
		o.start = v
		return nil
	}
}

func WithEndVersion(v int64) ListOptions {
	return func(o *listOptions) error {
		o.end = v
		return nil
	}
}

func WithDeleted() ListOptions {
	return func(o *listOptions) error {
		o.deleted = true
		return nil
	}
}

func ListForEnv(ctx context.Context, envID ids.ID, db *gorm.DB, options ...ListOptions) ([]models.FeatureToggleEnv, error) {
	opts := &listOptions{}
	for _, opt := range options {
		opt(opts)
	}
	query := db.WithContext(ctx).Select("feature_toggle_id, MAX(version) as version").Where("environment_id = ?", envID)
	if opts.start > 0 {
		query = query.Where("version > ?", opts.start)
	}
	if opts.end > 0 {
		query = query.Where("version <= ?", opts.end)
	}

	if opts.deleted {
		query = query.Unscoped()
	}

	type Pair struct {
		FeatureToggleID ids.ID
		Version         int64
	}
	var limitedFtEnvs []Pair
	if err := query.Group("feature_toggle_id").Table("feature_toggle_envs").Find(&limitedFtEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	if len(limitedFtEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	var ftEnvs []models.FeatureToggleEnv
	where := make([][]interface{}, len(limitedFtEnvs))
	for i, pair := range limitedFtEnvs {
		where[i] = []interface{}{pair.FeatureToggleID, pair.Version}
	}
	ftEnvQuery := db.WithContext(ctx)
	if opts.deleted {
		ftEnvQuery = ftEnvQuery.Unscoped()
	}
	if err := ftEnvQuery.Where("environment_id = ?", envID).Where("(feature_toggle_id, version) in ?", where).Order("created_at DESC").Preload("CreatedBy").Preload("FeatureToggle", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("FeatureToggle.CreatedBy").Find(&ftEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	return ftEnvs, nil
}

func GetHistoryForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) ([]models.FeatureToggleEnv, error) {
	var ftEnvs []models.FeatureToggleEnv
	if err := db.WithContext(ctx).Where("environment_id = ? AND feature_toggle_id = ?", envID, id).Order("created_at DESC").Preload("CreatedBy").Preload("FeatureToggle", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Find(&ftEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	if len(ftEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	return ftEnvs, nil
}

type PbOpts struct {
	FillUser bool
}

func platforms(ft *models.FeatureToggle) []pb_ft.Platform_Type {
	var platforms []pb_ft.Platform_Type
	if ft.IsServer {
		platforms = append(platforms, pb_ft.Platform_DEFAULT)
	}
	if ft.IsWeb {
		platforms = append(platforms, pb_ft.Platform_WEB)
	}
	if ft.IsIOS {
		platforms = append(platforms, pb_ft.Platform_IOS)
	}
	if ft.IsAndroid {
		platforms = append(platforms, pb_ft.Platform_ANDROID)
	}
	return platforms
}

func MultiPb(ctx context.Context, ftEnvs []models.FeatureToggleEnv, ory *ory.Ory, opts PbOpts) ([]*pb_ft.FeatureToggle, error) {
	var fts []*pb_ft.FeatureToggle
	for _, ftEnv := range ftEnvs {
		pb, err := Pb(ctx, &ftEnv, ory, PbOpts{FillUser: false})
		if err != nil {
			return nil, err
		}
		fts = append(fts, pb)
	}
	return fts, nil
}

func Pb(ctx context.Context, ftEnv *models.FeatureToggleEnv, ory *ory.Ory, opts PbOpts) (*pb_ft.FeatureToggle, error) {
	ft := ftEnv.FeatureToggle
	var pbCreatedBy, pbUpdatedBy *pb_user.User
	if opts.FillUser {
		if ft.CreatedBy.OryID != "" {
			createdByIdenty, err := users.FetchIdentity(ctx, ft.CreatedBy.OryID, ory.Api())
			if err != nil {
				return nil, err
			}
			pbCreatedBy = users.Pb(createdByIdenty, &ft.CreatedBy)
			// Let's filter out some fields
			users.LimitedPbUser(pbCreatedBy)
		}
		if ftEnv.CreatedBy.OryID != "" {
			updatedByIdenty, err := users.FetchIdentity(ctx, ftEnv.CreatedBy.OryID, ory.Api())
			if err != nil {
				return nil, err
			}

			pbUpdatedBy = users.Pb(updatedByIdenty, &ftEnv.CreatedBy)
			users.LimitedPbUser(pbUpdatedBy)
		}
	}

	if ftEnv.DeletedAt.Valid {
		return &pb_ft.FeatureToggle{
			Id: string(ft.ID),
			// Name is important because we use it as an index
			Name:      ft.Name,
			CreatedAt: timestamppb.New(ft.CreatedAt),
			UpdatedAt: timestamppb.New(ftEnv.CreatedAt),
			DeletedAt: timestamppb.New(ftEnv.DeletedAt.Time),
			Platforms: platforms(&ft),
		}, nil
	}
	res := &pb_ft.FeatureToggle{
		Id:          string(ft.ID),
		CreatedAt:   timestamppb.New(ft.CreatedAt),
		UpdatedAt:   timestamppb.New(ftEnv.CreatedAt),
		Name:        ft.Name,
		Description: ft.Description,
		ToggleType:  ft.Type,
		Version:     ftEnv.Version,
		Enabled:     ftEnv.Enabled,
		Platforms:   platforms(&ft),
		ProjectId:   string(ft.ProjectID),
		CreatedBy:   pbCreatedBy,
		UpdatedBy:   pbUpdatedBy,
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
