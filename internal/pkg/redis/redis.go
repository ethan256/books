package redis

import (
	"context"
	"time"

	"github.com/ethan256/books/configs"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func NewClient() (*redis.Client, error) {
	cfg := configs.Get().Redis
	return redisConnect(cfg.Addr, cfg.Pass, cfg.Db, cfg.MaxRetries, cfg.PoolSize, cfg.MinIdleConns)
}

func redisConnect(addr, pass string, db, maxRetries, poolSize, minIdleConns int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           db,
		MaxRetries:   maxRetries,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis err")
	}

	return client, nil
}
