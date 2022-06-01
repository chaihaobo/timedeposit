// Package model
// @author： Boice
// @createTime：2022/5/27 16:47
package dto

import "time"

type RetryFlowReqModel struct {
	FlowIdList []string `json:"flow_id_list"`
}

type RetryFlowSearchModel struct {
	Page   *Page `json:"page"`
	Search *Search
}

type Search struct {
	AccountId string `json:"account_id"`
}

func DefaultRetryFlowSearchModel() *RetryFlowSearchModel {
	return &RetryFlowSearchModel{Page: DefaultPage(), Search: new(Search)}
}

type FailFlowModel struct {
	Id              uint      `json:"id"`
	FlowId          string    `json:"flow_id"`
	AccountId       string    `json:"account_id"`
	FlowName        string    `json:"flow_name"`
	FlowStatus      string    `json:"flow_status"`
	FailedOperation string    `json:"failed_operation"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}
