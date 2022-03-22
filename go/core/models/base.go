package models

import (
	"time"

	id "stackv2/go/core/ids"

	"gorm.io/gorm"
)

type Model struct {
	ID        id.ID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
