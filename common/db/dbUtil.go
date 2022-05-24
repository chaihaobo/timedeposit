/*
 * @Author: Hugo
 * @Date: 2022-05-06 01:47:09
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:33:32
 */
package db

import (
	"fmt"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/util"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Build the db handle for database
var _db *gorm.DB
var dbOnce sync.Once

// Initial gorm
func initDB() {
	// Generate Mysql DSN
	// Mysql dsn format： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	// replace values which like {username}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.TDConf.Db.Username, config.TDConf.Db.Password, config.TDConf.Db.Host, config.TDConf.Db.Port, config.TDConf.Db.Db)
	var err error
	// Connect Mysql, Get DB connection
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		util.CheckAndExit(err)
	}

	sqlDB, _ := _db.DB()

	// Set the db connection configuration
	sqlDB.SetMaxOpenConns(100) // set the max openning connection number
	sqlDB.SetMaxIdleConns(20)  // set the max idle connection number
}

// Get db connection
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		initDB()
	})
	return _db
}
