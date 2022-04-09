package models

import (
	"stackv2/go/core/ids"

	"gorm.io/gorm"
)

type Environment struct {
	Model
	Name        string
	Description string
	ProjectID   string
	Project     Project
}

func (m Environment) ObjectType() ids.ObjectType {
	return ids.Environment
}

func (m Environment) BeforeCreate(tx *gorm.DB) error {
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AllModels = append(AllModels, &Environment{})
}
