package models

import (
	"platform/go/core/ids"
	"time"

	"gorm.io/gorm"
)

type ApiKey struct {
	Model
	Name          string
	ExpiresAt     time.Time
	ProjectID     ids.ID
	Project       Project
	EnvironmentID ids.ID `gorm:"index"`
	Environment   Environment
}

func (m ApiKey) ObjectType() ids.ObjectType {
	return ids.ApiKey
}

func (m ApiKey) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&ApiKey{})
}
