/*
 * @Author: Hugo
 * @Date: 2022-05-06 08:58:01
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-06 10:19:17
 */
package db

import "time"

type TFlowNode struct {
	ID         int64     `gorm:"column:id" db:"column:id" json:"id" form:"id"`
	FlowName   string    `gorm:"column:flow_name" db:"column:flow_name" json:"flow_name" form:"flow_name"`
	NodeName   string    `gorm:"column:node_name" db:"column:node_name" json:"node_name" form:"node_name"`
	NodePath   string    `gorm:"column:node_path" db:"column:node_path" json:"node_path" form:"node_path"`
	NodeDetail string    `gorm:"column:node_detail" db:"column:node_detail" json:"node_detail" form:"node_detail"`
	CreateTime time.Time `gorm:"column:create_time" db:"column:create_time" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" db:"column:update_time" json:"update_time" form:"update_time"`
}
