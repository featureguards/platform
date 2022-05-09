package models

import (
	"platform/go/core/ids"

	"gorm.io/gorm"
)

type Environment struct {
	Model
	Name        string
	Description string
	ProjectID   ids.ID
	Project     Project
}

func (m Environment) ObjectType() ids.ObjectType {
	return ids.Environment
}

func (m Environment) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&Environment{})
}
