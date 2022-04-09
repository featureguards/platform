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
	OwnerID        string
	Owner          User
	ProjectMembers []ProjectMember
}

func (m Project) ObjectType() ids.ObjectType {
	return ids.Project
}

func (m Project) BeforeCreate(tx *gorm.DB) error {
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type ProjectMember struct {
	Model
	ProjectID string `gorm:"index"`
	Project   Project
	UserID    string `gorm:"index"`
	User      User
	Role      pb_project.Project_Role
}

func (m ProjectMember) ObjectType() ids.ObjectType {
	return ids.ProjectMember
}

func (m ProjectMember) BeforeCreate(tx *gorm.DB) error {
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

type ProjectInvite struct {
	Model
	ProjectID string
	Project   Project
	Email     string `gorm:"index"`
	Status    pb_project.ProjectInvite_Status
}

func (m ProjectInvite) ObjectType() ids.ObjectType {
	return ids.ProjectInvite
}

func (m ProjectInvite) BeforeCreate(tx *gorm.DB) error {
	return beforeCreate(m.ID, m.ObjectType(), tx)
}

func init() {
	AllModels = append(AllModels, &ProjectInvite{}, &ProjectMember{}, &Project{})
}
