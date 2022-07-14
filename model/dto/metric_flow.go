// Package dto
// @author： Boice
// @createTime：2022/7/14 11:15
package dto

type FlowMetricQueryModel struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type FlowMetricResultModel struct {
	CreateDate string `json:"date"`
	TaskCnt    int    `json:"task_cnt"`
	SuccessCnt int    `json:"success_cnt"`
	RunningCnt int    `json:"running_cnt"`
	FailCnt    int    `json:"fail_cnt"`
	IsCall     int    `json:"is_call"`
}
