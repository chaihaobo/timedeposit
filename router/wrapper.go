// Package router
// @author： Boice
// @createTime：2022/7/14 10:14
package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/constant"
	"gitlab.com/bns-engineering/td/router/api"
	"net/http"
)

func GinWrapper(target func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		data, err := target(context)
		if err != nil {
			switch err.(type) {
			case constant.ServiceError:
				serviceError := err.(constant.ServiceError)
				context.JSON(serviceError.StatusCode, api.ErrorWithServiceError(serviceError, data))
				return
			default:
				context.JSON(http.StatusInternalServerError, api.Error(err.Error()))
				return
			}
		}
		context.JSON(http.StatusOK, api.SuccessData(data))

	}
}
