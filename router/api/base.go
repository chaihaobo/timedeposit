// Package api
// @author： Boice
// @createTime：
package api

const (
	SUCCESS_CODE = 200
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func success() *Response {
	return successData(nil)
}

func successData(data interface{}) *Response {
	return &Response{
		Code:    SUCCESS_CODE,
		Data:    data,
		Message: "success",
	}
}
