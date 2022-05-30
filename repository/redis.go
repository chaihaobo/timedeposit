// Package repository
// @author： Boice
// @createTime：2022/5/27 12:01
package repository

import (
	"context"
	"encoding/json"
	"gitlab.com/bns-engineering/td/common/cache"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	"go.uber.org/zap"
	"time"
)

var redisRepository = new(RedisRepository)

const (
	tdAccountPrefix      = "TD:ACCOUNT:"
	benefitAccountPrefix = "BENEFIT:ACCOUNT:"
)

type IRedisRepository interface {
	SaveTDAccount(account *mambuEntity.TDAccount) error
	SaveBenefitAccount(account *mambuEntity.TDAccount) error
	GetTDAccount(accountId string) *mambuEntity.TDAccount
	GetBenefitAccount(accountId string) *mambuEntity.TDAccount
}

type RedisRepository struct {
}

func (r *RedisRepository) SaveTDAccount(account *mambuEntity.TDAccount) error {
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	cache.GetRedis().Set(context.Background(), tdAccountPrefix+account.ID, string(marshal), time.Hour).Result()
	return nil
}

func (r *RedisRepository) SaveBenefitAccount(account *mambuEntity.TDAccount) error {
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	cache.GetRedis().Set(context.Background(), benefitAccountPrefix+account.ID, string(marshal), time.Hour)
	return nil
}

func (r *RedisRepository) GetTDAccount(accountId string) *mambuEntity.TDAccount {
	val := cache.GetRedis().Get(context.Background(), tdAccountPrefix+accountId).Val()
	account := new(mambuEntity.TDAccount)
	err := json.Unmarshal([]byte(val), account)
	if err != nil {
		zap.L().Error("get td account cache error ", zap.Error(err))
		return nil
	}
	return account
}

func (r *RedisRepository) GetBenefitAccount(accountId string) *mambuEntity.TDAccount {
	val := cache.GetRedis().Get(context.Background(), benefitAccountPrefix+accountId).Val()
	account := new(mambuEntity.TDAccount)
	err := json.Unmarshal([]byte(val), account)
	if err != nil {
		zap.L().Error("get td account cache error ", zap.Error(err))
		return nil
	}
	return account
}

func GetRedisRepository() IRedisRepository {
	return redisRepository
}
