package models

import (
	"platform/go/core/ids"

	"gorm.io/gorm"
)

// Sharded by OryID
type User struct {
	Model
	OryID string `gorm:"uniqueIndex"`
}

func (m User) ObjectType() ids.ObjectType {
	return ids.User
}

func (m User) BeforeCreate(tx *gorm.DB) error {
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&User{})
}
