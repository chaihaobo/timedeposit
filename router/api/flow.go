// Package api
// @author： Boice
// @createTime：2022/5/27 10:03
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/td/core/engine"
	mambuservices "gitlab.com/bns-engineering/td/service/mambuServices"
	"go.uber.org/zap"
	"net/http"
)

func StartFlow(c *gin.Context) {
	//Get all td accounts which need to process
	tmpQueryParam := generateSearchTDAccountParam()
	tmpTDAccountList, err := mambuservices.GetTDAccountListByQueryParam(tmpQueryParam)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Query mambu service for TD Account List failed! error: %v", err))
		return
	}
	if len(tmpTDAccountList) == 0 {
		zap.L().Info("Query mambu service for TD Account List get empty! No TD Account need to process")
		return
	}

	for _, tmpTDAcc := range tmpTDAccountList {
		zap.L().Info(fmt.Sprintf("Before Run Flow for Account: %v", tmpTDAcc.ID))
	}
	zap.L().Info("=======================================")

	for _, tmpTDAcc := range tmpTDAccountList {
		err := engine.Pool.Submit(func() {
			engine.Start(tmpTDAcc.ID)
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
		}
	}

}
