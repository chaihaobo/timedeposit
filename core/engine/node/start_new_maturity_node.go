// Package node
// @author： Boice
// @createTime：2022/5/26 17:59
package node

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/core/engine/mambu/accountservice"
	"gitlab.com/bns-engineering/td/core/engine/mambu/holidayservice"
	"strconv"
	"time"
)

const (
	endDayOfBigMonth   = 31
	endDayOfSmallMonth = 30
)

type StartNewMaturityNode struct {
	*Node
}

func (node *StartNewMaturityNode) Run(ctx context.Context) (INodeResult, error) {

	account, err := node.GetMambuAccount(ctx, node.AccountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseA() {
		activationDate := account.ActivationDate
		maturityDate, err := generateMaturityDateStr(node.GetContext(ctx), account.OtherInformation.Tenor, account.MaturityDate, account.MatureOnHoliday(), &activationDate)
		if err != nil {
			log.Info(ctx, fmt.Sprintf("Generate New Maturity Date failed, Account: %v", account.ID))
			return nil, err
		}
		note := fmt.Sprintf("TDE-AUTO-%v", node.FlowId)

		// create new maturity date
		_, err = accountservice.ChangeMaturityDate(node.GetContext(ctx), account.ID, maturityDate, note)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Update maturity date failed for account: %v", account.ID), err)
			return nil, errors.New("start New Maturity Date Failed")
		}
	} else {
		log.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}

// Calcuate the new maturity date by tenor
func generateMaturityDateStr(ctx context.Context, tenor string, maturityDate time.Time, matureOnHoliday bool, activationDate *time.Time) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error for convert tenor to int, tenor: %v", tenor), err)
		return "", errors.New("convert tenor to int failed")
	}
	resultDate := carbon.NewCarbon(maturityDate).AddMonths(tenorInt)
	// if last maturity day of month is 31 .then next maturity day is last of that month
	if activationDate != nil {
		resultDate.DaysInMonth()
		if carbon.NewCarbon(*activationDate).Day() == endDayOfBigMonth &&
			(resultDate.Day() == endDayOfSmallMonth || (resultDate.Month() == 2) && resultDate.Day() == resultDate.LastDayOfMonth().Day()) {
			resultDate = resultDate.LastDayOfMonth()
		}
	}
	holidayList := holidayservice.GetHolidayList(ctx)
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
