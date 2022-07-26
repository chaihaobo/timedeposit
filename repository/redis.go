// Package repository
// @author： Boice
// @createTime：2022/5/27 12:01
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common"
	"gitlab.com/bns-engineering/td/model/mambu"
	"go.uber.org/zap"
	"time"
)

const (
	tdAccountPrefix      = "TD:ACCOUNT:"
	benefitAccountPrefix = "BENEFIT:ACCOUNT:"
	idempotencyKeyPrefix = "idempotencyKey:"
	terminalRNNPrefix    = "terminalRNN:"
)

type IRedisRepository interface {
	SaveTDAccount(ctx context.Context, account *mambu.TDAccount, flowId string) error
	SaveBenefitAccount(ctx context.Context, account *mambu.TDAccount, flowId string) error
	GetTDAccount(ctx context.Context, flowId string) *mambu.TDAccount
	GetBenefitAccount(ctx context.Context, flowId string) *mambu.TDAccount
	GetIdempotencyKey(ctx context.Context, flowId string, nodeName string) string
	GetTerminalRRN(ctx context.Context, id string, name string, rrnGenerator func() string) string
}

type redisRepository struct {
	common *common.Common
}

func (r *redisRepository) SaveTDAccount(ctx context.Context, account *mambu.TDAccount, flowId string) error {
	tr := tracer.StartTrace(ctx, "redis_repository-SaveTDAccount")
	ctx = tr.Context()
	defer tr.Finish()
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	_, err = r.common.Cache.Set(context.Background(), tdAccountPrefix+flowId, string(marshal), time.Hour).Result()
	return err
}

func (r *redisRepository) SaveBenefitAccount(ctx context.Context, account *mambu.TDAccount, flowId string) error {
	tr := tracer.StartTrace(ctx, "redis_repository-SaveBenefitAccount")
	ctx = tr.Context()
	defer tr.Finish()
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	r.common.Cache.Set(context.Background(), benefitAccountPrefix+flowId, string(marshal), time.Hour)
	return nil
}

func (r *redisRepository) GetTDAccount(ctx context.Context, flowId string) *mambu.TDAccount {
	tr := tracer.StartTrace(ctx, "redis_repository-GetTDAccount")
	ctx = tr.Context()
	defer tr.Finish()
	val := r.common.Cache.Get(context.Background(), tdAccountPrefix+flowId).Val()
	if val == "" {
		return nil
	}
	account := new(mambu.TDAccount)
	err := json.Unmarshal([]byte(val), account)
	if err != nil {
		r.common.Logger.Info(ctx, "get td account cache error ", zap.Error(err))
		return nil
	}
	return account
}

func (r *redisRepository) GetIdempotencyKey(ctx context.Context, flowId string, nodeName string) string {
	tr := tracer.StartTrace(ctx, fmt.Sprintf("redis_repository-GetIdempotencyKey-%s", nodeName))
	ctx = tr.Context()
	defer tr.Finish()
	val := r.common.Cache.Get(ctx, idempotencyKeyPrefix+flowId+nodeName).Val()
	if val == "" {
		val = uuid.New().String()
		r.common.Cache.Set(ctx, idempotencyKeyPrefix+flowId+nodeName, val, time.Hour*24*30)
	}
	return val
}

func (r *redisRepository) GetTerminalRRN(ctx context.Context, flowId string, nodeName string, rrnGenerator func() string) string {
	tr := tracer.StartTrace(ctx, fmt.Sprintf("redis_repository-GetTerminalRRN-%s", nodeName))
	ctx = tr.Context()
	defer tr.Finish()
	val := r.common.Cache.Get(ctx, terminalRNNPrefix+flowId+nodeName).Val()
	if val == "" {
		val = rrnGenerator()
		r.common.Cache.Set(ctx, terminalRNNPrefix+flowId+nodeName, val, time.Hour*24*30)
	}
	return val
}

func (r *redisRepository) GetBenefitAccount(ctx context.Context, flowId string) *mambu.TDAccount {
	tr := tracer.StartTrace(ctx, "redis_repository-GetBenefitAccount")
	ctx = tr.Context()
	defer tr.Finish()
	val := r.common.Cache.Get(context.Background(), benefitAccountPrefix+flowId).Val()
	if val == "" {
		return nil
	}
	account := new(mambu.TDAccount)
	err := json.Unmarshal([]byte(val), account)
	if err != nil {
		r.common.Logger.Error(ctx, "get td account cache error ", err)
		return nil
	}
	return account
}

func newRedisRepository(common *common.Common) IRedisRepository {
	return &redisRepository{
		common: common,
	}
}
