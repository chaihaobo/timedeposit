/*
 * @Author: Hugo
 * @Date: 2022-05-17 11:39:14
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 03:45:14
 */
package po

import "time"

type TFlowAccountMaturityDay struct {
	Id           uint      `gorm:"column:id" db:"id" json:"id" form:"id"`
	FlowId       string    `gorm:"column:flow_id" db:"flow_id" json:"flow_id" form:"flow_id"`
	AccountId    string    `gorm:"column:account_id" db:"account_id" json:"account_id" form:"account_id"`
	MaturityDate time.Time `gorm:"column:maturity_date" db:"maturity_date" json:"maturity_date" form:"maturity_date"`
	CreateTime   time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
}
