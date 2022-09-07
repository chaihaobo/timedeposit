// Package middleware
// @author： Boice
// @createTime：2022/6/20 18:33
package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/common"
	"net/http"
)

const (
	authHeader = "Authorization"
)

func AuthMiddleware(common *common.Common) gin.HandlerFunc {
	return func(context *gin.Context) {
		authToken := context.GetHeader(authHeader)
		if "" == authToken || authToken != common.Credential.Server.AuthToken {
			context.Abort()
			context.String(http.StatusUnauthorized, "unauthorized")
		}
		context.Next()

	}
}
