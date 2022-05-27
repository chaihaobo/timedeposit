// Package model
// @author： Boice
// @createTime：2022/5/27 10:55
package model

import "time"

type TFlowNodeQueryLog struct {
	ID         int64     `gorm:"column:id" db:"column:id" json:"id" form:"id"`
	FLowId     string    `gorm:"column:flow_id" db:"column:flow_id" json:"flow_name" form:"flow_id"`
	NodeName   string    `gorm:"column:node_name" db:"column:node_name" json:"node_name" form:"node_name"`
	QueryType  string    `gorm:"column:query_type" db:"column:query_type" json:"query_type" form:"query_type"`
	Data       string    `gorm:"column:data" db:"column:data" json:"data" form:"data"`
	CreateTime time.Time `gorm:"column:create_time" db:"column:create_time" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" db:"column:update_time" json:"update_time" form:"update_time"`
}
