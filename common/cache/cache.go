// Package cache
// @author： Boice
// @createTime：2022/7/22 15:27
package cache

import (
	"github.com/go-redis/redis/v8"
	"gitlab.com/bns-engineering/td/common/config"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(config *config.Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	return &Cache{
		client,
	}
}
