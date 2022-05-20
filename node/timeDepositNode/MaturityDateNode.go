/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:43
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 05:04:04
 */
package timeDepositNode

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/bns-engineering/td/common/log"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/node"
	"gitlab.com/bns-engineering/td/service/mambuEntity"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
)

//Modify the Maturity Date for TD Account
type MaturityDateNode struct {
	node.Node
}

func (node *MaturityDateNode) Process() {
	CurNodeName := "maturity_date_node"
	tmpTDAccount, tmpFlowTask, nodeLog := node.GetAccAndFlowLog(CurNodeName)

	if !tmpTDAccount.IsCaseA() {
		log.Log.Info("No need to update maturity date, accNo: %v", tmpTDAccount.ID)
		node.UpdateLogWhenSkipNode(tmpFlowTask, CurNodeName, nodeLog)
	} else {
		var err error
		tmpTDAccount, err = updateMaturityDate(tmpTDAccount, tmpFlowTask.FlowId)
		if err != nil {
			log.Log.Error("Update maturity date failed for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFailed(tmpFlowTask, nodeLog, err)
			return
		} else {
			log.Log.Info("Finish update maturity date for account: %v", tmpTDAccount.ID)
			node.UpdateLogWhenNodeFinish(tmpFlowTask, nodeLog)
		}
	}
	log.Log.Debug("MaturityDateNode: OutputData: %v", tmpTDAccount)
	node.Node.Output <- tmpTDAccount
}

// Update maturity date for this account
func updateMaturityDate(tmpTDAccount mambuEntity.TDAccount, flowID string) (mambuEntity.TDAccount, error) {
	//Undo Maturity Date
	undoMaturityResult := mambuservices.UndoMaturityDate(tmpTDAccount.ID)
	if !undoMaturityResult {
		return tmpTDAccount, errors.New("undo Maturity Date Failed")
	}

	note := fmt.Sprintf("TDE-AUTO-%v", flowID)
	maturityDate, err := generateMaturityDateStr(tmpTDAccount.Otherinformation.Tenor, tmpTDAccount.Maturitydate)
	if err != nil {
		return tmpTDAccount, err
	}
	newTDAccount, err := mambuservices.ChangeMaturityDate(tmpTDAccount.ID, maturityDate, note)
	if err != nil {
		return tmpTDAccount, errors.New("start New Maturity Date Failed")
	}
	return newTDAccount, nil
}

//Calcuate the new maturity date by tenor
func generateMaturityDateStr(tenor string, maturityDate time.Time) (string, error) {
	tenorInt, err := strconv.Atoi(tenor)
	if err != nil {
		log.Log.Error("Error for convert tenor to int, tenor: %v", tenor)
		return "", errors.New("convert tenor to int failed")
	}
	resultDate := maturityDate.AddDate(0, tenorInt, 0)
	return util.GetDate(resultDate), nil
}
