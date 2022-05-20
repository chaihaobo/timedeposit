/*
 * @Author: Hugo
 * @Date: 2022-05-19 06:58:20
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 07:19:31
 */
package model

import "time"

type TFlowTransactions struct {
	Id                 uint      `gorm:"column:id" db:"id" json:"id" form:"id"`
	TransId            string    `gorm:"column:trans_id" db:"trans_id" json:"trans_id" form:"trans_id"`
	TerminalRrn        string    `gorm:"column:terminal_rrn" db:"terminal_rrn" json:"terminal_rrn" form:"terminal_rrn"`
	SourceAccountNo    string    `gorm:"column:source_account_no" db:"source_account_no" json:"source_account_no" form:"source_account_no"`
	SourceAccountName  string    `gorm:"column:source_account_name" db:"source_account_name" json:"source_account_name" form:"source_account_name"`
	BenefitAccountNo   string    `gorm:"column:benefit_account_no" db:"benefit_account_no" json:"benefit_account_no" form:"benefit_account_no"`
	BenefitAccountName string    `gorm:"column:benefit_account_name" db:"benefit_account_name" json:"benefit_account_name" form:"benefit_account_name"`
	Amount             float64   `gorm:"column:amount" db:"amount" json:"amount" form:"amount"`
	Channel            string    `gorm:"column:channel" db:"channel" json:"channel" form:"channel"`
	TransactionType    string    `gorm:"column:transaction_type" db:"transaction_type" json:"transaction_type" form:"transaction_type"`
	Result             int       `gorm:"column:result" db:"result" json:"result" form:"result"` //1:succeed, 0:failed
	EncodedKey         string    `gorm:"column:encoded_key" db:"encoded_key" json:"encoded_key" form:"encoded_key"`
	ErrorMsg           string    `gorm:"column:error_msg" db:"error_msg" json:"error_msg" form:"error_msg"`
	CreateTime         time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`
}
