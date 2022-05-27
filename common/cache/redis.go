// Package cache
// @author： Boice
// @createTime：2022/5/27 11:47

package cache

import (
	"github.com/go-redis/redis/v8"
	"gitlab.com/bns-engineering/td/common/config"
	"sync"
)

var (
	client     *redis.Client
	clientOnce sync.Once
)

func GetRedis() *redis.Client {
	clientOnce.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     config.TDConf.Redis.Addr,
			Password: config.TDConf.Redis.Password,
			DB:       config.TDConf.Redis.DB,
			PoolSize: config.TDConf.Redis.PoolSize,
		})
	})
	return client
}
