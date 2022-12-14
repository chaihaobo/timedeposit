// Package http
// @author： Boice
// @createTime：2022/7/26 14:24
package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/constant"
)

const (
	SuccessCode = 200
	ErrorCode   = 500
)

var OK interface{} = nil

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

func ErrorWithServiceError(serviceError constant.ServiceError, data interface{}) *Response {
	return &Response{
		Code:    serviceError.Code,
		Data:    data,
		Message: serviceError.Message,
	}

}

func SuccessData(data interface{}) *Response {
	return &Response{
		Code:    SuccessCode,
		Data:    data,
		Message: "success",
	}
}

type Controller interface {
	Apply(engine *gin.Engine)
}
