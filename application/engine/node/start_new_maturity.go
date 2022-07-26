// Package node
// @author： Boice
// @createTime：2022/5/26 17:59
package node

import (
	"context"
	"fmt"
	carbonv2 "github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/model/mambu"
	"strconv"
	"time"
)

const (
	genTensorMonthSize = 1000
)

type StartNewMaturityNode struct {
	node
}

func (node *StartNewMaturityNode) Run(ctx context.Context) (INodeResult, error) {

	account, err := node.GetMambuAccount(ctx, node.accountId, false)
	if err != nil {
		return nil, err
	}
	if account.IsCaseA(node.taskCreateTime) {
		activationDate := account.ActivationDate
		maturityDate, err := node.generateMaturityDateStr(node.GetContext(ctx), account.OtherInformation.Tenor, account.MaturityDate, account.MatureOnHoliday(), activationDate)
		if err != nil {
			node.common.Logger.Info(ctx, fmt.Sprintf("Generate New Maturity Date failed, Account: %v", account.ID))
			return nil, err
		}
		note := fmt.Sprintf("TDE-AUTO-%v", node.flowId)

		// create new maturity date
		_, err = node.service.Mambu.Account.ChangeMaturityDate(node.GetContext(ctx), account.ID, maturityDate, note)
		if err != nil {
			node.common.Logger.Error(ctx, fmt.Sprintf("Update maturity date failed for account: %v", account.ID), err)
			return nil, errors.New("start New Maturity Date Failed")
		}
	} else {
		node.common.Logger.Info(ctx, "not match! skip it")
		return ResultSkip, nil
	}
	return ResultSuccess, nil
}

// Calcuate the new maturity date by tenor
func (node *StartNewMaturityNode) generateMaturityDateStr(ctx context.Context, tenor string, maturityDate time.Time, matureOnHoliday bool, activationDate time.Time) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		node.common.Logger.Error(ctx, fmt.Sprintf("Error for convert tenor to int, tenor: %v", tenor), err)
		return "", errors.New("convert tenor to int failed")
	}
	resultDate := carbon.NewCarbon(getNextMaturityDay(activationDate, maturityDate, tenorInt))
	holidayList := node.service.Mambu.Holiday.GetHolidayList(ctx)
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

func generateTensorArray(tensorCount int, acticationTime time.Time) []time.Time {
	maturityTensorArray := make([]time.Time, 0)
	acticationDate := carbonv2.Parse(carbonv2.Time2Carbon(acticationTime).ToDateString())
	day := acticationDate.Day()
	for i := 0; i < tensorCount; i++ {
		acticationDate = acticationDate.SetDay(1)
		acticationDate = acticationDate.AddMonthsNoOverflow(1)
		if acticationDate.DaysInMonth() >= day {
			acticationDate = acticationDate.SetDay(day)
		} else {
			acticationDate = acticationDate.SetDay(acticationDate.DaysInMonth())
		}
		maturityTensorArray = append(maturityTensorArray, acticationDate.Carbon2Time())
	}
	return maturityTensorArray
}

func getNextMaturityDay(acticationDate, currentMaturityDate time.Time, tensor int) time.Time {
	tensorArray := generateTensorArray(genTensorMonthSize, acticationDate)
	for i, tensorTime := range tensorArray {
		if tensorTime.After(currentMaturityDate) {
			return tensorArray[i-1+tensor]
		}
	}
	return time.Now()
}

func isHoliday(time time.Time, holidayInfo mambu.HolidayInfo) bool {
	for _, holiday := range holidayInfo.Holidays {
		holidayDate := carbonv2.Parse(holiday.Date)
		carbonTime := carbonv2.Time2Carbon(time)
		if carbonTime.IsSameDay(holidayDate) {
			return true
		}
		if holiday.IsAnnuallyRecurring == true && carbonTime.Month() == holidayDate.Month() && carbonTime.DayOfMonth() == holidayDate.DayOfMonth() {
			return true
		}

	}
	return false
}
