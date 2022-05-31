package setup

import (
	"context"
	"platform/go/test/ory"

	kratos "github.com/ory/kratos-client-go"
)

func (a *Apps) CreateUser(ctx context.Context) (*kratos.Identity, error) {
	return ory.CreateIdentity(ctx, a.App.Ory().AdminApi())
}

func (a *Apps) CreateUserWithSession(ctx context.Context) (*kratos.Session, string, error) {
	return ory.CreateIdentityWithSession(ctx, a.App.Ory().Api())
}
