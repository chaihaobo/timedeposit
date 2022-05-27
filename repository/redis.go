// Package repository
// @author： Boice
// @createTime：2022/5/27 12:01
package repository

import (
	"context"
	"gitlab.com/bns-engineering/td/common/cache"
	"time"
)

var redisRepository = new(RedisRepository)

type IRedisRepository interface {
	Set(key string, value string)
	Get(key string) string
}

type RedisRepository struct {
}

func (r *RedisRepository) Set(key string, value string) {
	cache.GetRedis().Set(context.Background(), key, value, time.Hour*24*15)
}

func (r *RedisRepository) Get(key string) string {
	return cache.GetRedis().Get(context.Background(), key).String()
}

func GetRedisRepository() IRedisRepository {
	return redisRepository
}
