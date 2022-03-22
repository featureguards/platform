package models

type Project struct {
	Model
	Name         string
	Description  string
	Environments []Environment
	OwnerID      string
	Owner        User
}
