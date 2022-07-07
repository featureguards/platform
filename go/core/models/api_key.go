package models

import (
	"platform/go/core/ids"
	"time"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
	"github.com/lib/pq"

	"gorm.io/gorm"
)

func PbPlatformTypes(v pq.Int32Array) []pb_ft.Platform_Type {
	r := make([]pb_ft.Platform_Type, len(v))
	for i, e := range v {
		r[i] = pb_ft.Platform_Type(e)
	}
	return r
}

func PlatformTypes(v []pb_ft.Platform_Type) pq.Int32Array {
	r := make(pq.Int32Array, len(v))
	for i, e := range v {
		r[i] = int32(e)
	}
	return r
}

type ApiKey struct {
	Model
	Name          string
	ExpiresAt     time.Time
	ProjectID     ids.ID
	Project       Project
	EnvironmentID ids.ID `gorm:"index"`
	Environment   Environment
	Key           string        `gorm:"uniqueIndex"`
	Platforms     pq.Int32Array `gorm:"type:integer[]"`
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
