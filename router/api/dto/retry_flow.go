// Package model
// @author： Boice
// @createTime：2022/5/27 16:47
package dto

import "time"

type RetryFlowModel struct {
	FlowIdList []string `json:"flow_id_list"`
}

type FailFlowModel struct {
	Id          uint      `json:"id"`
	FlowId      string    `json:"flow_id"`
	AccountId   string    `json:"account_id"`
	CurNodeName string    `json:"cur_node_name"`
	CurStatus   string    `json:"cur_status"`
	CreateTime  time.Time `json:"create_time"`
}
