package models

type Environment struct {
	Model
	Name        string
	Description string
	ProjectID   string
	Project     Project
}
