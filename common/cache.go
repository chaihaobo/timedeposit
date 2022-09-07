// Package cache
// @author： Boice
// @createTime：2022/7/22 15:27
package common

import (
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	*redis.Client
}

func newCache(config *Config, credential *Credential) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     credential.Redis.Addr,
		Password: credential.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	return &Cache{
		client,
	}
}
