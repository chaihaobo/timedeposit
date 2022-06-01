// Package model
// @author： Boice
// @createTime：
package po

import "time"

type TFailTransactions struct {
	Id                 int64     `gorm:"column:id" db:"id" json:"id" form:"id"`
	FlowTransactionId  int64     `gorm:"column:flow_transaction_id" db:"flow_transaction_id" json:"flow_transaction_id" form:"flow_transaction_id"`
	TransId            string    `gorm:"column:trans_id" db:"trans_id" json:"trans_id" form:"trans_id"`
	TerminalRrn        string    `gorm:"column:terminal_rrn" db:"terminal_rrn" json:"terminal_rrn" form:"terminal_rrn"`
	SourceAccountNo    string    `gorm:"column:source_account_no" db:"source_account_no" json:"source_account_no" form:"source_account_no"`
	SourceAccountName  string    `gorm:"column:source_account_name" db:"source_account_name" json:"source_account_name" form:"source_account_name"`
	BenefitAccountNo   string    `gorm:"column:benefit_account_no" db:"benefit_account_no" json:"benefit_account_no" form:"benefit_account_no"`
	BenefitAccountName string    `gorm:"column:benefit_account_name" db:"benefit_account_name" json:"benefit_account_name" form:"benefit_account_name"`
	Amount             float64   `gorm:"column:amount" db:"amount" json:"amount" form:"amount"`
	Channel            string    `gorm:"column:channel" db:"channel" json:"channel" form:"channel"`
	TransactionType    string    `gorm:"column:transaction_type" db:"transaction_type" json:"transaction_type" form:"transaction_type"`
	RetryStatus        int       `gorm:"column:retry_status" db:"retry_status" json:"retry_status" form:"retry_status"`
	RetryTimes         int       `gorm:"column:retry_times" db:"retry_times" json:"retry_times" form:"retry_times"`
	CreateTime         time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
}
