// Package repository
// @author： Boice
// @createTime：2022/5/27 12:01
package repository

import (
	"context"
	"encoding/json"
	"gitlab.com/bns-engineering/td/common/cache"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"time"
)

var redisRepository = new(RedisRepository)

const (
	tdAccountPrefix      = "TD:ACCOUNT:"
	benefitAccountPrefix = "BENEFIT:ACCOUNT:"
)

type IRedisRepository interface {
	SaveTDAccount(account *mambu.TDAccount, flowId string) error
	SaveBenefitAccount(account *mambu.TDAccount, flowId string) error
	GetTDAccount(flowId string) *mambu.TDAccount
	GetBenefitAccount(flowId string) *mambu.TDAccount
}

type RedisRepository struct {
}

func (r *RedisRepository) SaveTDAccount(account *mambu.TDAccount, flowId string) error {
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	_, err = cache.GetRedis().Set(context.Background(), tdAccountPrefix+flowId, string(marshal), time.Hour).Result()
	return err
}

func (r *RedisRepository) SaveBenefitAccount(account *mambu.TDAccount, flowId string) error {
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	cache.GetRedis().Set(context.Background(), benefitAccountPrefix+flowId, string(marshal), time.Hour)
	return nil
}

func (r *RedisRepository) GetTDAccount(flowId string) *mambu.TDAccount {
	val := cache.GetRedis().Get(context.Background(), tdAccountPrefix+flowId).Val()
	if val == "" {
		return nil
	}
	account := new(mambu.TDAccount)
	err := json.Unmarshal([]byte(val), account)
	if err != nil {
		zap.L().Info("get td account cache error ", zap.Error(err))
		return nil
	}
	return account
}

func (r *RedisRepository) GetBenefitAccount(flowId string) *mambu.TDAccount {
	val := cache.GetRedis().Get(context.Background(), benefitAccountPrefix+flowId).Val()
	if val == "" {
		return nil
	}
	account := new(mambu.TDAccount)
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
