// Package api
// @author： Boice
// @createTime：
package api

const (
	SuccessCode = 200
	ErrorCode   = 500
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success() *Response {
	return SuccessData(nil)
}

func Error(message string) *Response {
	return &Response{
		Code:    ErrorCode,
		Message: message,
	}
}

func SuccessData(data interface{}) *Response {
	return &Response{
		Code:    SuccessCode,
		Data:    data,
		Message: "success",
	}
}
