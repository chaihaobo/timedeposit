/*
 * @Author: Hugo
 * @Date: 2022-05-07 02:49:23
 * @Last Modified by:   Hugo
 * @Last Modified time: 2022-05-07 02:49:23
 */
package model

import "time"

type TFlowNodeRelation struct {
	ID         int64     `gorm:"column:id" db:"column:id" json:"id" form:"id"`
	FlowName   string    `gorm:"column:flow_name" db:"column:flow_name" json:"flow_name" form:"flow_name"`
	NodeName   string    `gorm:"column:node_name" db:"column:node_name" json:"node_name" form:"node_name"`
	ResultCode string    `gorm:"column:result_code" db:"column:result_code" json:"result_code" form:"result_code"`
	NextNode   string    `gorm:"column:next_node" db:"column:next_node" json:"next_node" form:"next_node"`
	CreateTime time.Time `gorm:"column:create_time" db:"column:create_time" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" db:"column:update_time" json:"update_time" form:"update_time"`
}
