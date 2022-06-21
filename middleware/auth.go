// Package middleware
// @author： Boice
// @createTime：2022/6/20 18:33
package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common/config"
	"net/http"
)

const (
	authHeader = "Authorization"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authToken := context.GetHeader(authHeader)
		if "" == authToken || authToken != config.TDConf.Server.AuthToken {
			context.Abort()
			context.String(http.StatusUnauthorized, "unauthorized")
		}
		context.Next()

	}
}
