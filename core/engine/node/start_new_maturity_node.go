// Package node
// @author： Boice
// @createTime：2022/5/26 17:59
package node

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/bns-engineering/td/common/util"
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
		maturityDate, err := generateMaturityDateStr(node.GetContext(), account.OtherInformation.Tenor, account.MaturityDate)
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
	}
	return ResultSuccess, nil
}

// Calcuate the new maturity date by tenor
func generateMaturityDateStr(cxt context.Context, tenor string, maturityDate time.Time) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error for convert tenor to int, tenor: %v", tenor))
		return "", errors.New("convert tenor to int failed")
	}
	// Todo: note, should add logic for holidays
	resultDate := maturityDate.AddDate(0, tenorInt, 0)
	for _, tmpHoliday := range holidayservice.GetHolidayList(cxt) {
		if util.InSameDay(tmpHoliday, resultDate) {
			resultDate = maturityDate.AddDate(0, 0, 1)
		}
	}
	return util.GetDate(resultDate), nil
}
