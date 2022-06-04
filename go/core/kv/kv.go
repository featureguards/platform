package kv

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"platform/go/core/random"
	pb_private "platform/go/proto/private"
	pb_project "platform/go/proto/project"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type KeyType string

const (
	// Must all be unique
	ApiKey             KeyType = "api-key"
	RefreshToken       KeyType = "refresh-token"
	RefreshTokenFamily KeyType = "refresh-token-family"
	EnvironmentVersion KeyType = "env-version"
	EnvironmentToggles KeyType = "env-toggles"
	defaultExpiration          = time.Duration(time.Hour * 24)
	emptyDuration              = time.Duration(0)
)

var (
	ErrNotFound = errors.New("not found")
	exp         = 5 * time.Second
)

type KV struct {
	redis      *redis.Client
	expiration time.Duration
	rnd        *rand.Rand
	mu         sync.Mutex
}

type keyOpts struct {
	expiration time.Duration
	suffix     string
}

type Pending struct {
	key   string
	val   []byte
	redis *redis.Client
}

type KeyOptions func(*keyOpts)

func WithExpiration(exp time.Duration) KeyOptions {
	return func(ko *keyOpts) {
		ko.expiration = exp
	}
}

func WithSuffix(s string) KeyOptions {
	return func(ko *keyOpts) {
		ko.suffix = s
	}
}

type Opts struct {
	Redis      *redis.Client
	Expiration time.Duration
}

func New(opts Opts) (*KV, error) {
	exp := opts.Expiration
	if exp == emptyDuration {
		exp = defaultExpiration
	}
	return &KV{redis: opts.Redis, expiration: exp, rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}, nil
}

func (kv *KV) redisKey(keyType KeyType, k, suffix string) string {
	sb := strings.Builder{}
	sb.Grow(len(k) + len(keyType) + 8)
	sb.WriteRune('{')
	sb.Write([]byte(keyType))
	sb.WriteString("::")
	sb.WriteString(k)
	sb.WriteRune('}')
	if suffix != "" {
		sb.WriteString(suffix)
	}
	return sb.String()
}

func (kv *KV) pendingKey(keyType KeyType, k, suffix string) string {
	return kv.redisKey(keyType, k, suffix) + "-pending"
}

func (kv *KV) keyExp(opts *keyOpts) time.Duration {
	if opts.expiration != emptyDuration {
		return opts.expiration
	}
	return kv.expiration
}

func (kv *KV) SetNX(ctx context.Context, keyType KeyType, k string, v []byte, options ...KeyOptions) (bool, error) {
	opts := &keyOpts{}
	for _, opt := range options {
		opt(opts)
	}
	set, err := kv.redis.SetNX(ctx, kv.redisKey(keyType, k, opts.suffix), v, kv.keyExp(opts)).Result()
	if err != nil {
		return false, errors.WithStack(err)
	}
	return set, nil
}

func (kv *KV) SetProto(ctx context.Context, keyType KeyType, k string, m proto.Message, options ...KeyOptions) error {
	v, err := proto.Marshal(m)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := kv.SetNX(ctx, keyType, k, v, options...); err != nil {
		return err
	}
	return nil
}

func (kv *KV) Get(ctx context.Context, keyType KeyType, k string, options ...KeyOptions) ([]byte, error) {
	opts := &keyOpts{}
	for _, opt := range options {
		opt(opts)
	}
	res, err := kv.redis.MGet(ctx, kv.redisKey(keyType, k, opts.suffix), kv.pendingKey(keyType, k, opts.suffix)).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// If we have a pending key, return no elements
	if res[1] != nil {
		return nil, ErrNotFound
	}
	// Key not found
	if res[0] == nil {
		return nil, ErrNotFound
	}
	return []byte(res[0].(string)), nil
}

func (kv *KV) GetProto(ctx context.Context, keyType KeyType, k string, options ...KeyOptions) (proto.Message, error) {
	var m proto.Message
	switch keyType {
	case ApiKey:
		m = &pb_project.ApiKey{}
	case EnvironmentVersion:
		m = &pb_private.EnvironmentVersion{}
	case EnvironmentToggles:
		m = &pb_private.EnvironmentFeatureToggles{}
	default:
		return nil, errors.WithStack(fmt.Errorf("unknown key-type: %s", keyType))
	}
	v, err := kv.Get(ctx, keyType, k, options...)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(v, m); err != nil {
		return nil, errors.WithStack(err)
	}
	return m, nil
}

func (kv *KV) random() ([]byte, error) {
	kv.mu.Lock()
	data := random.RandBytes(16, kv.rnd)
	kv.mu.Unlock()
	return data, nil
}

func (kv *KV) StartPending(ctx context.Context, kt KeyType, k string, options ...KeyOptions) (*Pending, error) {
	opts := &keyOpts{}
	for _, opt := range options {
		opt(opts)
	}
	v, err := kv.random()
	if err != nil {
		return nil, err
	}
	pendingKey := kv.pendingKey(kt, k, opts.suffix)
	set, err := kv.redis.SetNX(ctx, pendingKey, v, exp).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !set {
		return nil, errors.WithStack(errors.New("concurrent start pending"))
	}

	// Clear the value key now that we got the lock
	if err := kv.redis.Del(ctx, kv.redisKey(kt, k, opts.suffix)).Err(); err != nil {
		return nil, errors.WithStack(errors.New("couldn't delete key"))
	}

	return &Pending{redis: kv.redis, key: pendingKey, val: v}, nil
}

func (p *Pending) Finish(ctx context.Context) error {
	if err := delKeyConditional.Run(ctx, p.redis, []string{p.key}, p.val).Err(); err != nil {
		log.Errorf("%s\n", err)
		return errors.WithStack(err)
	}
	return nil
}

var delKeyConditional = redis.NewScript(`
	if redis.call("get",KEYS[1]) == ARGV[1] then
		return redis.call("del",KEYS[1])
	else
		return 0
	end
	`)
