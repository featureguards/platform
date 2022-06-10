package stubs

import (
	"context"
	"platform/go/grpc/auth"
	"platform/go/grpc/middleware/jwt_auth"
	"platform/go/grpc/middleware/meta"
	"platform/go/grpc/middleware/token_auth"
	"platform/go/test/setup"
	"testing"

	pb_project "platform/go/proto/project"

	pb_user "github.com/featureguards/featureguards-go/v2/proto/user"
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
	return &Stubs{App: app}
}

func (s *Stubs) Create(ctx context.Context) error {
	if err := s.createUser(ctx); err != nil {
		return err
	}
	if err := s.createProject(ctx); err != nil {
		return err
	}
	if err := s.createApiKey(ctx); err != nil {
		return err
	}
	if err := s.CreateFeatureToggle(ctx); err != nil {
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
	md = md.Set(auth.ApiKeyMD, s.ApiKey.Key)
	return md.ToOutgoing(ctx)
}

func (s *Stubs) WithJwtToken(ctx context.Context, token string) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set(jwt_auth.Key, "Bearer "+token)
	return md.ToOutgoing(ctx)
}
