package models

import (
	"stackv2/go/core/ids"
	pb_project "stackv2/go/proto/project"

	"gorm.io/gorm"
)

type Project struct {
	Model
	Name           string
	Description    string
	Environments   []Environment
	OwnerID        ids.ID
	Owner          User
	ProjectMembers []ProjectMember
}

func (m Project) ObjectType() ids.ObjectType {
	return ids.Project
}

func (m Project) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.OwnerID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type ProjectMember struct {
	Model
	ProjectID ids.ID `gorm:"index"`
	Project   Project
	UserID    ids.ID `gorm:"index"`
	User      User
	Role      pb_project.Project_Role
}

func (m ProjectMember) ObjectType() ids.ObjectType {
	return ids.ProjectMember
}

func (m ProjectMember) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID, m.UserID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type ProjectInvite struct {
	Model
	ProjectID ids.ID
	Project   Project
	Email     string `gorm:"index"`
	Status    pb_project.ProjectInvite_Status
}

func (m ProjectInvite) ObjectType() ids.ObjectType {
	return ids.ProjectInvite
}

func (m ProjectInvite) BeforeCreate(tx *gorm.DB) error {
	var toCheck []ids.ID = []ids.ID{m.ID, m.ProjectID}
	for _, id := range toCheck {
		if err := id.Validate(); err != nil {
			return err
		}
	}
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AddModel(&ProjectInvite{})
	AddModel(&ProjectMember{})
	AddModel(&Project{})
}
