package app_context

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwt"
	kratos "github.com/ory/kratos-client-go"
)

type sessionKey struct{}
type jwtTokenKey struct{}

func WithSession(ctx context.Context, session *kratos.Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}
func SessionFromContext(ctx context.Context) (*kratos.Session, bool) {
	session, ok := ctx.Value(sessionKey{}).(*kratos.Session)
	return session, ok
}

func WithJwtToken(ctx context.Context, token jwt.Token) context.Context {
	return context.WithValue(ctx, jwtTokenKey{}, token)
}

func JwtTokenFromContext(ctx context.Context) (jwt.Token, bool) {
	token, ok := ctx.Value(jwtTokenKey{}).(jwt.Token)
	return token, ok
}
