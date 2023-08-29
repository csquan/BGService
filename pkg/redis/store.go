package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"user/pkg/util"
)

type redisCli interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}
type Store struct {
	cli    redisCli
	prefix string
}

func NewStore(prefix string) *Store {
	return &Store{
		Cli,
		prefix,
	}
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}

// Get ...
func (s *Store) Get(ctx context.Context, key string) (string, util.Err) {
	cmd := s.cli.Get(ctx, s.wrapperKey(key))
	result, err := cmd.Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return result, nil
		}
		return result, util.ErrRDB
	}
	return result, nil
}

// Set ...
func (s *Store) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) util.Err {
	cmd := s.cli.Set(ctx, s.wrapperKey(key), value, expiration)
	if cmd.Err() != nil {
		return util.ErrRDB
	}
	return nil
}

// Expire ...
func (s *Store) Expire(ctx context.Context, key string, expiration time.Duration) util.Err {
	cmd := s.cli.Expire(ctx, s.wrapperKey(key), expiration)
	if cmd.Err() != nil {
		return util.ErrRDB
	}
	return nil
}

// Delete ...
func (s *Store) Delete(ctx context.Context, key string) (bool, util.Err) {
	cmd := s.cli.Del(ctx, s.wrapperKey(key))
	if cmd.Err() != nil {
		return false, util.ErrRDB
	}
	return cmd.Val() > 0, nil
}

// Exists ...
func (s *Store) Exists(ctx context.Context, key string) (bool, util.Err) {
	cmd := s.cli.Exists(ctx, s.wrapperKey(key))
	if cmd.Err() != nil {
		return false, util.ErrRDB
	}
	return cmd.Val() > 0, nil
}

// Close ...
func (s *Store) Close() error {
	return s.cli.Close()
}
