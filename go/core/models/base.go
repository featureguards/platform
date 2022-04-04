package models

import (
	"errors"
	"time"

	"stackv2/go/core/ids"

	"google.golang.org/protobuf/reflect/protoreflect"
	"gorm.io/gorm"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrNoSession         = errors.New("no session")
	ErrNoID              = errors.New("no id")
	ErrInvalidObjectType = errors.New("invalid object type")
	ErrInvalidID         = errors.New("invlid id")
)

type ModelObject interface {
	ObjectType() ids.ObjectType
}

type Model struct {
	ModelObject `gorm:"-"`
	ID          ids.ID `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		return ErrNoID
	}

	ot, _, err := ids.Parse(m.ID)
	if err != nil {
		return ErrInvalidID
	}
	modelOT := m.ObjectType()
	if ot.Validate() != nil || modelOT.Validate() != nil || ot != modelOT {
		return ErrInvalidObjectType
	}
	return nil
}

func FieldsFromPb(m protoreflect.Message) map[string]interface{} {
	fields := make(map[string]interface{})
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if v.IsValid() {
			fields[string(fd.Name())] = v.Interface()
		}
		return true
	})
	return fields
}
