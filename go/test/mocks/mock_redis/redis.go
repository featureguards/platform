package mock_redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func New(t *testing.T) (*redis.Client, *miniredis.Miniredis, error) {
	s := miniredis.NewMiniRedis()
	if err := s.Start(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	t.Cleanup(s.Close)
	c := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return c, s, nil
}
