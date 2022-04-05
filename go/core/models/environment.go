package models

import "stackv2/go/core/ids"

var (
	_ ModelObject = Environment{}
)

type Environment struct {
	Model
	Name        string
	Description string
	ProjectID   string
	Project     Project
}

func (m Environment) ObjectType() ids.ObjectType {
	return ids.Environment
}

func init() {
	AllModels = append(AllModels, &Environment{})
}
