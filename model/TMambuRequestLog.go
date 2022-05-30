// Package model
// @author： Boice
// @createTime：
package model

import "time"

type TMambuRequestLog struct {
	ID           int64     `gorm:"column:id" db:"column:id" json:"id" form:"id"`
	FlowId       string    `gorm:"column:flow_id" db:"column:flow_id" json:"flow_id" form:"flow_id"`
	NodeName     string    `gorm:"column:node_name" db:"column:node_name" json:"node_name" form:"node_name"`
	AccountId    string    `gorm:"column:account_id" db:"column:account_id" json:"account_id" form:"account_id"`
	Type         string    `gorm:"column:type" db:"column:type" json:"type" form:"type"`
	RequestUrl   string    `gorm:"column:request_url" db:"column:request_url" json:"request_url" form:"request_url"`
	RequestBody  string    `gorm:"column:request_body" db:"column:request_body" json:"request_body" form:"request_body"`
	ResponseCode int       `gorm:"column:response_code" db:"column:response_code" json:"response_code" form:"response_code"`
	ResponseBody string    `gorm:"column:response_body" db:"column:response_body" json:"response_body" form:"response_body"`
	Error        string    `gorm:"column:error" db:"column:error" json:"error" form:"error"`
	CreateTime   time.Time `gorm:"column:create_time" db:"column:create_time" json:"create_time" form:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" db:"column:update_time" json:"update_time" form:"update_time"`
}

// TMambuRequestLogsBuilder TMamboRequestLogs builder pattern code
type TMambuRequestLogsBuilder struct {
	tMambuRequestLogs *TMambuRequestLog
}

func NewTMambuRequestLogsBuilder() *TMambuRequestLogsBuilder {
	tMambuRequestLogs := &TMambuRequestLog{}
	b := &TMambuRequestLogsBuilder{tMambuRequestLogs: tMambuRequestLogs}
	return b
}

func (b *TMambuRequestLogsBuilder) ID(iD int64) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.ID = iD
	return b
}

func (b *TMambuRequestLogsBuilder) FlowId(flowId string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.FlowId = flowId
	return b
}

func (b *TMambuRequestLogsBuilder) NodeName(nodeName string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.NodeName = nodeName
	return b
}

func (b *TMambuRequestLogsBuilder) Type(reqType string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.Type = reqType
	return b
}

func (b *TMambuRequestLogsBuilder) RequestUrl(requestUrl string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.RequestUrl = requestUrl
	return b
}

func (b *TMambuRequestLogsBuilder) RequestBody(requestBody string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.RequestBody = requestBody
	return b
}

func (b *TMambuRequestLogsBuilder) ResponseCode(responseCode int) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.ResponseCode = responseCode
	return b
}

func (b *TMambuRequestLogsBuilder) ResponseBody(responseBody string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.ResponseBody = responseBody
	return b
}

func (b *TMambuRequestLogsBuilder) Error(error string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.Error = error
	return b
}

func (b *TMambuRequestLogsBuilder) CreateTime(createTime time.Time) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.CreateTime = createTime
	return b
}

func (b *TMambuRequestLogsBuilder) UpdateTime(updateTime time.Time) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.UpdateTime = updateTime
	return b
}

func (b *TMambuRequestLogsBuilder) AccountId(accountId string) *TMambuRequestLogsBuilder {
	b.tMambuRequestLogs.AccountId = accountId
	return b
}

func (b *TMambuRequestLogsBuilder) Build() *TMambuRequestLog {
	return b.tMambuRequestLogs
}
