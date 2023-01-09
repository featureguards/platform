package dynamic_settings

import (
	"context"
	"database/sql"
	"fmt"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/models/users"
	"platform/go/core/ory"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"
	pb_platform "github.com/featureguards/featureguards-go/v2/proto/platform"
	pb_user "github.com/featureguards/featureguards-go/v2/proto/user"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type GetDSOpts struct {
	WithUser bool
}

func Get(ctx context.Context, id ids.ID, db *gorm.DB, opts GetDSOpts) (*models.DynamicSetting, error) {
	var ds models.DynamicSetting
	query := db.WithContext(ctx).Where("id = ?", id)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.First(&ds).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return &ds, nil
}

func GetByName(ctx context.Context, projectID ids.ID, name string, db *gorm.DB, opts GetDSOpts) (*models.DynamicSetting, error) {
	var ds models.DynamicSetting
	query := db.WithContext(ctx).Where("project_id = ? AND name = ?", projectID, name)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.First(&ds).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return &ds, nil
}

func GetForEnv(ctx context.Context, envID, projectID ids.ID, db *gorm.DB, opts GetDSOpts) ([]models.DynamicSetting, error) {
	var dss []models.DynamicSetting
	query := db.WithContext(ctx).Where("environment_id = ? AND project_id = ?", envID, projectID)
	if opts.WithUser {
		query = query.Preload("CreatedBy")
	}
	if err := query.Find(&dss).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)

	}
	return dss, nil
}

func GetLatestForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) (*models.DynamicSettingEnv, error) {
	var dsEnv models.DynamicSettingEnv
	if err := db.WithContext(ctx).Where("dynamic_setting_id = ? AND environment_id = ?", id, envID).Order("version desc").Preload("CreatedBy").Preload("DynamicSetting").First(&dsEnv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &dsEnv, nil
}

func MaxVersionForEnv(ctx context.Context, envID ids.ID, db *gorm.DB) (int64, error) {
	var version sql.NullInt64
	if err := db.WithContext(ctx).Select("MAX(version) as version").Where("environment_id = ?", envID).Table("dynamic_setting_envs").Find(&version).Error; err != nil || !version.Valid {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, models.ErrNotFound
		}
		return -1, errors.WithStack(err)
	}
	return version.Int64, nil
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

func ListForEnv(ctx context.Context, envID ids.ID, db *gorm.DB, options ...ListOptions) ([]models.DynamicSettingEnv, error) {
	opts := &listOptions{}
	for _, opt := range options {
		opt(opts)
	}
	query := db.WithContext(ctx).Select("dynamic_setting_id, MAX(version) as version").Where("environment_id = ?", envID)
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
		DynamicSettingID ids.ID
		Version          int64
	}
	var limitedFtEnvs []Pair
	if err := query.Group("dynamic_setting_id").Table("dynamic_setting_envs").Find(&limitedFtEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	if len(limitedFtEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	var dsEnvs []models.DynamicSettingEnv
	where := make([][]interface{}, len(limitedFtEnvs))
	for i, pair := range limitedFtEnvs {
		where[i] = []interface{}{pair.DynamicSettingID, pair.Version}
	}
	dsEnvQuery := db.WithContext(ctx)
	if opts.deleted {
		dsEnvQuery = dsEnvQuery.Unscoped()
	}
	if err := dsEnvQuery.Where("environment_id = ?", envID).Where("(dynamic_setting_id, version) in ?", where).Order("created_at DESC").Preload("CreatedBy").Preload("DynamicSetting", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("DynamicSetting.CreatedBy").Find(&dsEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	return dsEnvs, nil
}

func GetHistoryForEnv(ctx context.Context, id, envID ids.ID, db *gorm.DB) ([]models.DynamicSettingEnv, error) {
	var dsEnvs []models.DynamicSettingEnv
	if err := db.WithContext(ctx).Where("environment_id = ? AND dynamic_setting_id = ?", envID, id).Order("created_at DESC").Preload("CreatedBy").Preload("DynamicSetting", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Find(&dsEnvs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	if len(dsEnvs) <= 0 {
		return nil, models.ErrNotFound
	}
	return dsEnvs, nil
}

type PbOpts struct {
	FillUser bool
}

func platforms(ds *models.DynamicSetting) []pb_platform.Type {
	var platforms []pb_platform.Type
	if ds.IsServer {
		platforms = append(platforms, pb_platform.Type_DEFAULT)
	}
	if ds.IsWeb {
		platforms = append(platforms, pb_platform.Type_WEB)
	}
	if ds.IsIOS {
		platforms = append(platforms, pb_platform.Type_IOS)
	}
	if ds.IsAndroid {
		platforms = append(platforms, pb_platform.Type_ANDROID)
	}
	return platforms
}

func MultiPb(ctx context.Context, dsEnvs []models.DynamicSettingEnv, ory *ory.Ory, opts PbOpts) ([]*pb_ds.DynamicSetting, error) {
	var dss []*pb_ds.DynamicSetting
	for _, dsEnv := range dsEnvs {
		pb, err := Pb(ctx, &dsEnv, ory, PbOpts{FillUser: false})
		if err != nil {
			return nil, err
		}
		dss = append(dss, pb)
	}
	return dss, nil
}

func Pb(ctx context.Context, dsEnv *models.DynamicSettingEnv, ory *ory.Ory, opts PbOpts) (*pb_ds.DynamicSetting, error) {
	ds := dsEnv.DynamicSetting
	var pbCreatedBy, pbUpdatedBy *pb_user.User
	if opts.FillUser {
		if ds.CreatedBy.OryID != "" {
			createdByIdenty, err := users.FetchIdentity(ctx, ds.CreatedBy.OryID, ory.Api())
			if err != nil {
				return nil, err
			}
			pbCreatedBy = users.Pb(createdByIdenty, &ds.CreatedBy)
			// Let's filter out some fields
			users.LimitedPbUser(pbCreatedBy)
		}
		if dsEnv.CreatedBy.OryID != "" {
			updatedByIdenty, err := users.FetchIdentity(ctx, dsEnv.CreatedBy.OryID, ory.Api())
			if err != nil {
				return nil, err
			}

			pbUpdatedBy = users.Pb(updatedByIdenty, &dsEnv.CreatedBy)
			users.LimitedPbUser(pbUpdatedBy)
		}
	}

	if dsEnv.DeletedAt.Valid {
		return &pb_ds.DynamicSetting{
			Id: string(ds.ID),
			// Name is important because we use it as an index
			Name:      ds.Name,
			CreatedAt: timestamppb.New(ds.CreatedAt),
			UpdatedAt: timestamppb.New(dsEnv.CreatedAt),
			DeletedAt: timestamppb.New(dsEnv.DeletedAt.Time),
			Platforms: platforms(&ds),
		}, nil
	}
	res := &pb_ds.DynamicSetting{
		Id:          string(ds.ID),
		CreatedAt:   timestamppb.New(ds.CreatedAt),
		UpdatedAt:   timestamppb.New(dsEnv.CreatedAt),
		Name:        ds.Name,
		Description: ds.Description,
		SettingType: ds.Type,
		Version:     dsEnv.Version,
		Platforms:   platforms(&ds),
		ProjectId:   string(ds.ProjectID),
		CreatedBy:   pbCreatedBy,
		UpdatedBy:   pbUpdatedBy,
	}
	pb, err := dsEnv.ProtoMessage(ds.Type)
	if err != nil {
		return nil, err
	}
	if err := fillProto(res, pb); err != nil {
		return nil, err
	}
	return res, nil
}

func fillProto(ds *pb_ds.DynamicSetting, pb proto.Message) error {
	switch ds.SettingType {
	case pb_ds.DynamicSetting_BOOL:
		ds.SettingDefinition = &pb_ds.DynamicSetting_BoolValue{BoolValue: pb.(*pb_ds.BoolValue)}
	case pb_ds.DynamicSetting_INTEGER:
		ds.SettingDefinition = &pb_ds.DynamicSetting_IntegerValue{IntegerValue: pb.(*pb_ds.IntegerValue)}
	case pb_ds.DynamicSetting_FLOAT:
		ds.SettingDefinition = &pb_ds.DynamicSetting_FloatValue{FloatValue: pb.(*pb_ds.FloatValue)}
	case pb_ds.DynamicSetting_STRING:
		ds.SettingDefinition = &pb_ds.DynamicSetting_StringValue{StringValue: pb.(*pb_ds.StringValue)}
	case pb_ds.DynamicSetting_LIST:
		ds.SettingDefinition = &pb_ds.DynamicSetting_ListValues{ListValues: pb.(*pb_ds.ListValues)}
	case pb_ds.DynamicSetting_SET:
		ds.SettingDefinition = &pb_ds.DynamicSetting_SetValues{SetValues: pb.(*pb_ds.SetValues)}
	case pb_ds.DynamicSetting_MAP:
		ds.SettingDefinition = &pb_ds.DynamicSetting_MapValues{MapValues: pb.(*pb_ds.MapValues)}
	default:
		return errors.WithStack(fmt.Errorf("unknown type %v", ds.SettingType))
	}
	return nil
}
