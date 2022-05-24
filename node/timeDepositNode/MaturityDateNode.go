/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:43
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:11:46
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"

	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Modify the Maturity Date for TD Account
type MaturityDateNode struct {
	node.Node
	// nodeName string
}

func NewMaturityDateNode() *MaturityDateNode {
	tmpNode := new(MaturityDateNode)
	// tmpNode.nodeName = "maturity_date_node"
	tmpNode.Node.NodeRun = tmpNode
	return tmpNode
}

func (node *MaturityDateNode) Process() {
	node.RunNode("maturity_date_node")
}

// Update maturity date for this account
func (node *MaturityDateNode) RunProcess(tmpTDAccount mambuEntity.TDAccount, flowID string, nodeName string) (constant.FlowNodeStatus, error) {
	// Skip updating maturity date or not?
	if !tmpTDAccount.IsCaseA() {
		zap.L().Info(fmt.Sprintf("No need to update maturity date, accNo: %v", tmpTDAccount.ID))
		return constant.FlowNodeSkip, nil
	}

	//Undo Maturity Date
	undoMaturityResult := mambuservices.UndoMaturityDate(tmpTDAccount.ID)
	if !undoMaturityResult {
		zap.L().Info(fmt.Sprintf("Undo Maturity Date for account failed: %v", tmpTDAccount.ID))
		return constant.FlowNodeFailed, errors.New("undo Maturity Date Failed")
	}

	//generate new Maturity Date
	maturityDate, err := generateMaturityDateStr(tmpTDAccount.OtherInformation.Tenor, tmpTDAccount.MaturityDate)
	if err != nil {
		zap.L().Info(fmt.Sprintf("Generate New Maturity Date failed, Account: %v", tmpTDAccount.ID))
		return constant.FlowNodeFailed, err
	}
	note := fmt.Sprintf("TDE-AUTO-%v", flowID)

	//create new maturity date
	_, err = mambuservices.ChangeMaturityDate(tmpTDAccount.ID, maturityDate, note)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Update maturity date failed for account: %v", tmpTDAccount.ID))
		return constant.FlowNodeFailed, errors.New("start New Maturity Date Failed")
	}
	return constant.FlowNodeFinish, nil
}

//Calcuate the new maturity date by tenor
func generateMaturityDateStr(tenor string, maturityDate time.Time) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error for convert tenor to int, tenor: %v", tenor))
		return "", errors.New("convert tenor to int failed")
	}
	//Todo: note, should add logic for holidays
	resultDate := maturityDate.AddDate(0, tenorInt, 0)
	for _, tmpHoliday := range mambuservices.GetHolidayList() {
		if util.InSameDay(tmpHoliday, resultDate) {
			resultDate = maturityDate.AddDate(0, 0, 1)
		}
	}
	return util.GetDate(resultDate), nil
}
