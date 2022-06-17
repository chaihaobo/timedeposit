// Package node
// @author： Boice
// @createTime：2022/5/26 17:59
package node

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/mambu/holidayservice"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type StartNewMaturityNode struct {
	*Node
}

func (node *StartNewMaturityNode) Run() (INodeResult, error) {

	account, err := node.GetMambuAccount(node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseA() {
		// generate new Maturity Date
		maturityDate, err := generateMaturityDateStr(node.GetContext(), account.OtherInformation.Tenor, account.MaturityDate, account.MatureOnHoliday())
		if err != nil {
			zap.L().Info(fmt.Sprintf("Generate New Maturity Date failed, Account: %v", account.ID))
			return nil, err
		}
		note := fmt.Sprintf("TDE-AUTO-%v", node.FlowId)

		// create new maturity date
		_, err = accountservice.ChangeMaturityDate(node.GetContext(), account.ID, maturityDate, note)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Update maturity date failed for account: %v", account.ID))
			return nil, errors.New("start New Maturity Date Failed")
		}
	} else {
		zap.L().Info("not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}

// Calcuate the new maturity date by tenor
func generateMaturityDateStr(cxt context.Context, tenor string, maturityDate time.Time, matureOnHoliday bool) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error for convert tenor to int, tenor: %v", tenor))
		return "", errors.New("convert tenor to int failed")
	}
	resultDate := carbon.NewCarbon(maturityDate).AddMonths(tenorInt)
	holidayList := holidayservice.GetHolidayList(cxt)
	if !matureOnHoliday {
		for {
			if resultDate.IsWeekend() || isHoliday(resultDate.Time, holidayList) {
				resultDate = resultDate.AddDays(1)
			} else {
				break
			}
		}

	}

	return resultDate.DateString(), nil
}

func isHoliday(time time.Time, holidays []time.Time) bool {
	for _, holiday := range holidays {
		if carbon.NewCarbon(time).IsSameDay(carbon.NewCarbon(holiday)) {
			return true
		}
	}
	return false
}
