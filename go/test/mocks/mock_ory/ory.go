package mock_ory

import (
	"platform/go/core/ory"
	"testing"
)

func New(t *testing.T) (*ory.Ory, error) {
	return ory.New(ory.Opts{})
}
