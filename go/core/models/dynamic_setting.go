package models

import (
	"fmt"
	"platform/go/core/ids"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type DynamicSetting struct {
	Model
	Description string
	ProjectID   ids.ID `gorm:"index:idx_dynamic_settings_project_id_name,unique"`
	Project     Project
	Name        string `gorm:"index:idx_dynamic_settings_project_id_name,unique"`
	CreatedByID ids.ID
	CreatedBy   User
	Type        pb_ds.DynamicSetting_Type // Immutable

	// PlatformType. Done this way to enable indexing in the future.
	IsIOS     bool // Immutable
	IsAndroid bool // Immutable
	IsWeb     bool // Immutable
	IsServer  bool // Immutable
}

func (m DynamicSetting) ObjectType() ids.ObjectType {
	return ids.DynamicSetting
}

func (m DynamicSetting) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID, m.CreatedByID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type DynamicSettingEnv struct {
	Model
	ProjectID        ids.ID
	Project          Project
	DynamicSettingID ids.ID `gorm:"index"`
	DynamicSetting   DynamicSetting
	EnvironmentID    ids.ID `gorm:"index"`
	Environment      Environment
	Version          int64 `gorm:"index"`
	CreatedByID      ids.ID
	CreatedBy        User
	Enabled          bool
	// One of BoolValue | StringValue | IntegerValue | FloatValue | SetValues | Map Values | ListValues
	Proto []byte
}

func (m DynamicSettingEnv) ProtoMessage(t pb_ds.DynamicSetting_Type) (proto.Message, error) {
	switch t {
	case pb_ds.DynamicSetting_BOOL:
		var obj pb_ds.BoolValue
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_FLOAT:
		var obj pb_ds.FloatValue
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_INTEGER:
		var obj pb_ds.IntegerValue
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_STRING:
		var obj pb_ds.StringValue
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_SET:
		var obj pb_ds.SetValues
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_LIST:
		var obj pb_ds.ListValues
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	case pb_ds.DynamicSetting_MAP:
		var obj pb_ds.MapValues
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			return nil, errors.WithStack(err)
		}
		return &obj, nil
	}
	err := errors.WithStack(fmt.Errorf("unknown feature flag type: %s", t))
	return nil, err
}

func (m DynamicSettingEnv) ObjectType() ids.ObjectType {
	return ids.DynamicSettingEnv
}

func (m DynamicSettingEnv) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID, m.DynamicSettingID, m.EnvironmentID, m.CreatedByID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&DynamicSetting{})
	AddModel(&DynamicSettingEnv{})
}
