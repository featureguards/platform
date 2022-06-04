package kv

import (
	"bytes"
	"context"

	pb_project "platform/go/proto/project"
	"platform/go/test/mocks/mock_redis"

	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestStartFinishPending(t *testing.T) {
	c, s, err := mock_redis.New(t)
	defer s.Close()
	assert.Nil(t, err, "%s", err)
	ctx := context.Background()
	cache, err := New(Opts{Redis: c})
	assert.Nil(t, err, "%s", err)

	apiKey := "123"
	pending, err := cache.StartPending(ctx, ApiKey, apiKey)
	assert.Nil(t, err)

	// Another StartPending must fail
	if _, err := cache.StartPending(ctx, ApiKey, apiKey); err == nil {
		t.Error(err)
	}

	v, err := c.Get(ctx, cache.pendingKey(ApiKey, apiKey, "")).Result()
	assert.Nil(t, err)

	if !bytes.Equal(pending.val, []byte(v)) {
		t.Errorf("Wrong pending value %v, wants %v\n", []byte(v), pending.val)
	}

	if err := pending.Finish(ctx); err != nil {
		t.Error(err)
	}
	_, err = c.Get(ctx, cache.pendingKey(ApiKey, apiKey, "")).Result()
	assert.EqualError(t, err, redis.Nil.Error())
}

func TestGetProto(t *testing.T) {
	c, s, err := mock_redis.New(t)
	defer s.Close()
	assert.Nil(t, err, "%s", err)
	ctx := context.Background()
	cache, err := New(Opts{Redis: c})
	assert.Nil(t, err)

	pb := &pb_project.ApiKey{
		Id: "1234",
	}

	found, err := cache.GetProto(ctx, ApiKey, pb.Id)
	if err != ErrNotFound {
		t.Errorf("Got %v and %v want %v and %v", found, err, nil, ErrNotFound)
	}

	// Set some value
	if err := cache.SetProto(ctx, ApiKey, pb.Id, pb); err != nil {
		t.Error(err)
	}

	found, err = cache.GetProto(ctx, ApiKey, pb.Id)
	if err != nil || !proto.Equal(found, pb) {
		t.Errorf("Got (%v, %v) want (%v, %v)", found, err, pb, nil)
	}

	pending, err := cache.StartPending(ctx, ApiKey, pb.Id)
	if err != nil {
		t.Error(err)
	}

	found, err = cache.GetProto(ctx, ApiKey, pb.Id)
	if err != ErrNotFound {
		t.Errorf("Got %v and %v want %v and %v", found, err, nil, ErrNotFound)
	}

	if err := pending.Finish(ctx); err != nil {
		t.Error(err)
	}

	_, err = cache.GetProto(ctx, ApiKey, pb.Id)
	assert.EqualError(t, err, ErrNotFound.Error())
}
