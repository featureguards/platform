package models

import (
	"errors"
	"log"
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

	AllModels = make(map[ids.ObjectType]ModelObject)
)

type ModelObject interface {
	ObjectType() ids.ObjectType
	BeforeCreate(tx *gorm.DB) error
}

type Model struct {
	ID        ids.ID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func beforeCreate(id ids.ID, modelOT ids.ObjectType, tx *gorm.DB) error {
	if id == "" {
		return ErrNoID
	}

	ot, _, err := ids.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
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

func AddModel(model ModelObject) {
	if _, ok := AllModels[model.ObjectType()]; ok {
		log.Fatalf("Duplicate models with ot=%s\n", model.ObjectType())
	}
	AllModels[model.ObjectType()] = model
}
