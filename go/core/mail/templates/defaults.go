package templates

import "os"

type defaults struct {
}

func (d defaults) FromEmail() string {
	return os.Getenv("FROM_EMAIL")
}

func (d defaults) FromName() string {
	return "FeatureGuards Team"
}
