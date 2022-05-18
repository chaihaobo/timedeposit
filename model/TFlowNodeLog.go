/*
 * @Author: Hugo
 * @Date: 2022-05-17 11:39:28
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:37:24
 */
package model

import "time"

type TFlowNodeLog struct {
	Id         uint      `gorm:"column:id" db:"id" json:"id" form:"id"`
	AccountId  string    `gorm:"column:account_id" db:"account_id" json:"account_id" form:"account_id"`
	FlowId     string    `gorm:"column:flow_id" db:"flow_id" json:"flow_id" form:"flow_id"`
	FlowName   string    `gorm:"column:flow_name" db:"flow_name" json:"flow_name" form:"flow_name"`
	NodeName   string    `gorm:"column:node_name" db:"node_name" json:"node_name" form:"node_name"`
	NodeResult string    `gorm:"column:node_result" db:"node_result" json:"node_result" form:"node_result"`
	NodeMsg    string    `gorm:"column:node_msg" db:"node_msg" json:"node_msg" form:"node_msg"`
	CreateTime time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
}
