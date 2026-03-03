package cacheadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

// Internal interface abstraction for the cache client used by RedisAdapter.
type (
	cacheClient interface {
		Exists(ctx context.Context, key string) bool
		Get(ctx context.Context, key string, value interface{}) error
		Set(item *cache.Item) error
		Delete(ctx context.Context, key string) error
	}

	// RedisAdapter is a small wrapper that namespaces and manages cache operations.
	RedisAdapter struct {
		ctx       context.Context
		client    cacheClient
		namespace string
	}

	// RedisConfig contains configuration for initializing RedisAdapter.
	RedisConfig struct {
		// DSN format: redis://<user>:<password>@<host>:<port>/<db_number>
		DSN       string
		Namespace string
	}
)

// NewRedis creates a new RedisAdapter using the provided configuration.
func NewRedis(c *RedisConfig) (*RedisAdapter, error) {
	if c.Namespace == "" {
		return nil, errors.New("namespace is required")
	}

	options, err := redis.ParseURL(c.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	redisClient := redis.NewClient(options)
	redisCache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return &RedisAdapter{
		ctx:       context.Background(),
		namespace: c.Namespace,
		client:    redisCache,
	}, nil
}

// Has reports whether the key exists in the cache under the configured namespace.
func (r *RedisAdapter) Has(key string) bool {
	return r.client.Exists(r.ctx, makeID(r.namespace, key))
}

// Save stores a value under the key with the specified TTL.
func (r *RedisAdapter) Save(key string, val interface{}, ttl time.Duration) error {
	return r.client.Set(&cache.Item{
		Ctx:   r.ctx,
		Key:   makeID(r.namespace, key),
		Value: val,
		TTL:   ttl,
	})
}

// Delete removes the key from the cache.
func (r *RedisAdapter) Delete(key string) error {
	return r.client.Delete(r.ctx, makeID(r.namespace, key))
}
