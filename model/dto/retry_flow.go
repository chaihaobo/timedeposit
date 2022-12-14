// Package model
// @author： Boice
// @createTime：2022/5/27 16:47
package dto

import (
	"time"
)

type RetryFlowReqModel struct {
	FlowIdList []string `json:"flow_id_list"`
}

type RetryFlowSearchModel struct {
	*Pagination
	*Search
}

type Search struct {
	AccountId string `json:"account_id" form:"account_id"`
}

func DefaultRetryFlowSearchModel() *RetryFlowSearchModel {
	return &RetryFlowSearchModel{
		DefaultPage(),
		&Search{},
	}
}

type FailFlowModel struct {
	Id              uint      `json:"id"`
	FlowId          string    `json:"flow_id"`
	AccountId       string    `json:"account_id"`
	FlowName        string    `json:"flow_name"`
	FlowStatus      string    `json:"flow_status"`
	FailedOperation string    `json:"failed_operation"`
	AmountToMove    string    `json:"amount_to_move"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}

type RetryResponseDTO struct {
	FlowCount     int                     `json:"flow_count"`
	FlowFailCount int                     `json:"flow_fail_count"`
	FailInfo      []RetryFailInfoResponse `json:"fail_info"`
}

type RetryFailInfoResponse struct {
	FlowId  string `json:"flow_id"`
	Message string `json:"message"`
}
