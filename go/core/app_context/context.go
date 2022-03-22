package app_context

import (
	"context"

	kratos "github.com/ory/kratos-client-go"
)

type sessionKey struct{}

func WithSession(ctx context.Context, session *kratos.Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}
func SessionFromContext(ctx context.Context) (*kratos.Session, bool) {
	session, ok := ctx.Value(sessionKey{}).(*kratos.Session)
	return session, ok
}
