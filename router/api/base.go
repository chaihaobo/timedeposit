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

func success() *Response {
	return successData(nil)
}

func error(message string) *Response {
	return &Response{
		Code:    ErrorCode,
		Message: message,
	}
}

func successData(data interface{}) *Response {
	return &Response{
		Code:    SuccessCode,
		Data:    data,
		Message: "success",
	}
}
