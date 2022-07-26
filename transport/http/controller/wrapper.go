// Package controller
// @author： Boice
// @createTime：2022/7/14 10:14
package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/constant"
	"net/http"
)

func ControlleWrapper(target func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		data, err := target(context)
		if err != nil {
			switch err.(type) {
			case constant.ServiceError:
				serviceError := err.(constant.ServiceError)
				context.JSON(serviceError.StatusCode, ErrorWithServiceError(serviceError, data))
				return
			default:
				context.JSON(http.StatusInternalServerError, Error(err.Error()))
				return
			}
		}
		context.JSON(http.StatusOK, SuccessData(data))

	}
}
