// Package repository
// @author： Boice
// @createTime：2022/5/26 13:41
package repository

import (
	"context"
	"gitlab.com/bns-engineering/common/tracer"
	"gitlab.com/bns-engineering/td/common/db"
	"gitlab.com/bns-engineering/td/model/po"
)

var flowAcountMaturityDayRepository *FlowAcountMaturityDayRepository

type IFlowAcountMaturityDayRepository interface {
	Save(ctx context.Context, flowAccountMaturity *po.TFlowAccountMaturityDay)
	GetFirstMaturityDate(ctx context.Context, accountId string, flowId string) *po.TFlowAccountMaturityDay
}

type FlowAcountMaturityDayRepository struct{}

func (f *FlowAcountMaturityDayRepository) Save(ctx context.Context, flowAccountMaturity *po.TFlowAccountMaturityDay) {
	tr := tracer.StartTrace(ctx, "flow_account_maturity_day-Save")
	ctx = tr.Context()
	defer tr.Finish()
	db.GetDB().Save(flowAccountMaturity)
}

func (f *FlowAcountMaturityDayRepository) GetFirstMaturityDate(ctx context.Context, accountId string, flowId string) *po.TFlowAccountMaturityDay {
	tr := tracer.StartTrace(ctx, "flow_account_maturity_day-GetLastMaturityDate")
	ctx = tr.Context()
	defer tr.Finish()
	last := new(po.TFlowAccountMaturityDay)
	db.GetDB().Where("account_id=? and flow_id!=?", accountId, flowId).First(last)
	if last.Id > 0 {
		return last
	}
	return nil
}

func init() {
	flowAcountMaturityDayRepository = new(FlowAcountMaturityDayRepository)
}

func GetFlowAcountMaturityDayRepository() IFlowAcountMaturityDayRepository {
	return flowAcountMaturityDayRepository
}
