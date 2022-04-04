package models

import "stackv2/go/core/ids"

var (
	_ ModelObject = User{}
)

// Sharded by OryID
type User struct {
	Model
	OryID string `gorm:"uniqueIndex"`
}

func (m User) ObjectType() ids.ObjectType {
	return ids.User
}
