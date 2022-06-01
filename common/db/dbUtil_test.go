/*
 * @Author: Hugo
 * @Date: 2022-05-06 01:47:09
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-06 10:10:05
 */
package db

import (
	"fmt"
	"gitlab.com/bns-engineering/td/model/db"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestGetDB(t *testing.T) {
	dbConn := GetDB()
	// testCreate(dbConn)
	testSelect(dbConn)
}

func testSelect(dbConn *gorm.DB) {
	var tmpNode db.TFlowNode
	dbConn.Where("flow_name = ?", "time_deposit_flow").First(&tmpNode)
	fmt.Println(tmpNode)
}

func testCreate(dbConn *gorm.DB) {
	tdEndNode := db.TFlowNode{
		FlowName:   "time_deposit_flow",
		NodeName:   "end_node",
		NodeDetail: "the final node of td flow, update",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	dbConn.Save(tdEndNode)
}
