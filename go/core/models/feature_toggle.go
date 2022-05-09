package models

import (
	"fmt"
	"platform/go/core/ids"
	pb_ft "platform/go/proto/feature_toggle"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type FeatureToggle struct {
	Model
	Description string
	ProjectID   ids.ID `gorm:"index:idx_feature_toggles_project_id_name,unique"`
	Project     Project
	Name        string `gorm:"index:idx_feature_toggles_project_id_name,unique"`
	CreatedByID ids.ID
	CreatedBy   User
	Type        pb_ft.FeatureToggle_Type // Immutable
	IsMobile    bool
	IsWeb       bool
}

func (m FeatureToggle) ObjectType() ids.ObjectType {
	return ids.FeatureToggle
}

func (m FeatureToggle) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID, m.CreatedByID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type FeatureToggleEnv struct {
	Model
	ProjectID       ids.ID
	Project         Project
	FeatureToggleID ids.ID `gorm:"index"`
	FeatureToggle   FeatureToggle
	EnvironmentID   ids.ID `gorm:"index"`
	Environment     Environment
	Version         int64 `gorm:"index"`
	CreatedByID     ids.ID
	CreatedBy       User
	Enabled         bool
	// One of OnOffFeature | PercentageFeature | PermissionFeature | ExperimentFeature
	Proto []byte
}

func (m FeatureToggleEnv) ProtoMessage(t pb_ft.FeatureToggle_Type) (proto.Message, error) {
	switch t {
	case pb_ft.FeatureToggle_EXPERIMENT:
		var obj pb_ft.ExperimentFeature
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			log.Error(errors.WithStack(err))
			return nil, err
		}
		return &obj, nil
	case pb_ft.FeatureToggle_ON_OFF:
		var obj pb_ft.OnOffFeature
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			log.Error(errors.WithStack(err))
			return nil, err
		}
		return &obj, nil
	case pb_ft.FeatureToggle_PERCENTAGE:
		var obj pb_ft.PercentageFeature
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			log.Error(errors.WithStack(err))
			return nil, err
		}
		return &obj, nil
	case pb_ft.FeatureToggle_PERMISSION:
		var obj pb_ft.PermissionFeature
		if err := proto.Unmarshal(m.Proto, &obj); err != nil {
			log.Error(errors.WithStack(err))
			return nil, err
		}
		return &obj, nil
	}
	err := errors.WithStack(fmt.Errorf("unknown feature toggle type: %s", t))
	log.Error(err)
	return nil, err
}

func (m FeatureToggleEnv) ObjectType() ids.ObjectType {
	return ids.FeatureToggleEnv
}

func (m FeatureToggleEnv) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID, m.FeatureToggleID, m.EnvironmentID, m.CreatedByID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&FeatureToggle{})
	AddModel(&FeatureToggleEnv{})
}
