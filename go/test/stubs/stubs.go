package stubs

import (
	"context"
	"platform/go/grpc/auth"
	"platform/go/grpc/middleware/meta"
	"platform/go/grpc/middleware/token_auth"
	"platform/go/test/setup"
	"testing"

	pb_project "platform/go/proto/project"

	pb_user "github.com/featureguards/client-go/proto/user"
	"github.com/stretchr/testify/require"
)

type Stubs struct {
	App *setup.Apps

	Token  string
	User   *pb_user.User
	Proj   *pb_project.Project
	ApiKey *pb_project.ApiKey
}

func New(ctx context.Context, t *testing.T) *Stubs {
	app := setup.App(t)
	stubs := &Stubs{App: app}
	err := stubs.create(ctx)
	require.Nil(t, err)
	return stubs
}

func (s *Stubs) create(ctx context.Context) error {
	if err := s.createUser(ctx); err != nil {
		return err
	}
	if err := s.createProject(ctx); err != nil {
		return err
	}
	if err := s.createApiKey(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Stubs) WithToken(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set(token_auth.Key, s.Token)
	return md.ToOutgoing(ctx)
}

func (s *Stubs) WithAPiKey(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set(auth.ApiKeyMD, s.ApiKey.Id)
	return md.ToOutgoing(ctx)
}
