package web_auth

import (
	"context"
	"net/http"

	"stackv2/go/core/app_context"
	"stackv2/go/core/ory"
	"stackv2/go/grpc/error_codes"
	"stackv2/go/grpc/middleware/meta"

	kratos "github.com/ory/kratos-client-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cookieKey  = "cookie"
	sessionKey = "app.sid"
)

type Auth struct {
	client *kratos.APIClient
}

type AuthOpts struct {
	KratosPublicURL string
}

func New(opts AuthOpts) (*Auth, error) {
	client, err := ory.New(ory.Opts{KratosPublicURL: opts.KratosPublicURL})
	if err != nil {
		return nil, err
	}
	return &Auth{client: client}, nil
}

func (a *Auth) Authenticate(ctx context.Context) (context.Context, error) {
	// TODO: Add partial support for paths based on x-envoy-original-path
	logger := log.WithContext(ctx)
	m := meta.ExtractIncoming(ctx)
	cookie := m.Get(cookieKey)
	if cookie == "" {
		return nil, status.Error(codes.Unauthenticated, "No cookie")
	}

	session, res, err := a.client.V0alpha2Api.ToSession(ctx).Cookie(cookie).Execute()
	if err != nil {
		logger.Warnf("Error Kratos session: %s", err)
		code := error_codes.GrpcCode(int32(res.StatusCode))
		return nil, status.Error(code, http.StatusText(res.StatusCode))
	}
	logger.Infof("%+v", session)
	if session.GetActive() {
		return app_context.WithSession(ctx, session), nil
	}
	return nil, status.Error(codes.Unauthenticated, "invalid session")
}
