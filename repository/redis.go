package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	rds *redis.Client
	ctx context.Context
}

func NewRedisRepository(rds *redis.Client, ctx context.Context) Repository {
	return &redisRepository{
		rds: rds,
		ctx: ctx,
	}
}

func (r *redisRepository) Save(str string, value any, exp int) error {
	return r.rds.SetEX(r.ctx, str, value, time.Duration(exp)*time.Second).Err()
}

func (r *redisRepository) Get(str string) any {
	return r.rds.Get(r.ctx, str).Val()
}
