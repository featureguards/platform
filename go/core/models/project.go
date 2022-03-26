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

// Sharded by Key
type ProjectInvite struct {
	Model
	ProjectID string
	Project   Project
	Email     string
	Status    int
	Key       string `gorm:"uniqueIndex"`
}
