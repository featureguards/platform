package scopes

import (
	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
	jwtx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
)

const (
	PlatformsScope = "platforms"
)

var (
	ErrNotFound = errors.New("not found")
)

func PlatformsArray(platforms []pb_ft.Platform_Type) []interface{} {
	res := make([]interface{}, len(platforms))
	for i, p := range platforms {
		res[i] = p
	}
	return res

}

func Platforms(scopesKey string, token jwtx.Token) (map[pb_ft.Platform_Type]struct{}, error) {
	v, found := token.Get(scopesKey)
	if !found {
		return nil, ErrNotFound
	}

	// It's a map[string]interface{}
	scopes, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.WithStack(errors.New("unexpected type for scopes"))
	}

	pp := scopes[PlatformsScope]
	if pp == nil {
		return nil, ErrNotFound
	}

	platforms, ok := pp.([]interface{})
	if !ok {
		return nil, errors.WithStack(errors.New("unexpected type for scopes"))
	}
	res := make(map[pb_ft.Platform_Type]struct{}, len(platforms))
	for _, p := range platforms {
		res[pb_ft.Platform_Type(p.(float64))] = struct{}{}
	}
	return res, nil
}
