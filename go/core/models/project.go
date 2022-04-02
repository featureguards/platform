package models

type Project struct {
	Model
	Name         string
	Description  string
	Environments []Environment
	OwnerID      string
	Owner        User
}

type ProjectMember struct {
	Model
	ProjectID string
	Project   Project
	UserID    string
	User      string
}

type ProjectInvite struct {
	Model
	ProjectID string
	Project   Project
	Email     string `gorm:"index"`
	Status    int
	Key       string `gorm:"uniqueIndex"`
}
