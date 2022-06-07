/*
 * @Author: Hugo
 * @Date: 2022-05-17 11:39:14
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 03:45:14
 */
package po

import "time"

type TFlowTaskInfo struct {
	Id          uint      `gorm:"column:id" db:"id" json:"id" form:"id"`
	FlowId      string    `gorm:"column:flow_id" db:"flow_id" json:"flow_id" form:"flow_id"`
	AccountId   string    `gorm:"column:account_id" db:"account_id" json:"account_id" form:"account_id"`
	FlowName    string    `gorm:"column:flow_name" db:"flow_name" json:"flow_name" form:"flow_name"`
	FlowStatus  string    `gorm:"column:flow_status" db:"flow_status" json:"flow_status" form:"flow_status"`
	CurNodeName string    `gorm:"column:cur_node_name" db:"cur_node_name" json:"cur_node_name" form:"cur_node_name"`
	CurStatus   string    `gorm:"column:cur_status" db:"cur_status" json:"cur_status" form:"cur_status"`
	EndStatus   string    `gorm:"column:end_status" db:"end_status" json:"end_status" form:"end_status"`
	StartTime   time.Time `gorm:"column:start_time" db:"start_time" json:"start_time" form:"start_time"`
	EndTime     time.Time `gorm:"column:end_time" db:"end_time" json:"end_time" form:"end_time"`
	Enable      bool      `gorm:"column:enable" db:"enable" json:"enable" form:"enable"`
	CreateTime  time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
}
