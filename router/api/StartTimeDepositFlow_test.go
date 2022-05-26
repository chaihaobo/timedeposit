/*
 * @Author: Hugo
 * @Date: 2022-05-16 09:08:29
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 11:25:01
 */
 package api

 import (
	 "testing"
 
	 "github.com/gin-gonic/gin"
	 commonConfig "gitlab.com/bns-engineering/td/common/config"
	 logger "gitlab.com/bns-engineering/td/common/log"
	 "gitlab.com/bns-engineering/td/flow"
	 "gitlab.com/bns-engineering/td/service/mambuEntity"
	 "go.uber.org/zap"
 )
 
 func init() {
	 conf := commonConfig.Setup("../../config.yaml")
	 logger.SetUp(conf)
	 zap.L().Info("===============Start Test Whole Flow==============")
	 flow.InitWorkflow()
 }
 
 
 func TestStartTDFlow(t *testing.T) {
	 type args struct {
		 c *gin.Context
	 }
	 tests := []struct {
		 name string
		 args args
	 }{
		 // TODO: Add test cases.
	 }
	 for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			 StartTDFlow(tt.args.c)
		 })
	 }
 }
 
 func TestRunFlow(t *testing.T) {
	 type args struct {
		 tmpTDAcc *mambuEntity.TDAccount
	 }
	 tests := []struct {
		 name string
		 args args
	 }{
		 // TODO: Add test cases.
		 {
			 name : "test Account : 11249460359, test Update Maturity Date",
			 args : args{ 
				 tmpTDAcc : &mambuEntity.TDAccount{
					 ID: "11249460359",
				 },
			 },
		 },
		 {
			 name : "test Account : 11747126703",
			 args : args{ 
				 tmpTDAcc : &mambuEntity.TDAccount{
					 ID: "11747126703",
				 },
			 },
		 },
		 {
			 name : "test Account : 11714744288",
			 args : args{ 
				 tmpTDAcc : &mambuEntity.TDAccount{
					 ID: "11714744288",
				 },
			 },
		 },
		 {
			 name : "test Account : 11842046986",
			 args : args{ 
				 tmpTDAcc : &mambuEntity.TDAccount{
					 ID: "11842046986",
				 },
			 },
		 },
	 }
	 for _, tt := range tests {
		 t.Run(tt.name, func(t *testing.T) {
			 RunFlow(tt.args.tmpTDAcc)
		 })
	 }
	 t.Log("run flow by account success")
 }
 
 