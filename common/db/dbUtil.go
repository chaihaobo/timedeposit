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
	"go.uber.org/zap"
	"sync"
	"time"

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
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		util.CheckAndExit(err)
	}

	sqlDB, _ := _db.DB()
	// Set the db connection configuration
	sqlDB.SetMaxOpenConns(config.TDConf.Db.MaxOpenConn) // set the max openning connection number
	sqlDB.SetMaxIdleConns(config.TDConf.Db.MaxIdleConn) // set the max idle connection number
	sqlDB.SetConnMaxLifetime(time.Hour)

}

// Get db connection
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		initDB()
	})
	return _db
}

func Paginate(pageNo int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo <= 0 {
			pageNo = 1
		}
		if pageSize > 100 {
			pageNo = 100
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func FindPage(db *gorm.DB, pageNo int, pageSize int, resultBind interface{}, totalBind *int64) {
	zap.L().Info("page query start")
	startTime := time.Now()
	countSql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {

		return tx.Count(totalBind)
	})
	querySql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(Paginate(pageNo, pageSize)).Find(resultBind)
	})

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		zap.L().Info("query total task start...")
		totalTaskStartTime := time.Now()
		// var count string
		defer wait.Done()
		GetDB().Raw(countSql).Scan(totalBind)
		useTime := time.Now().Sub(totalTaskStartTime).Milliseconds()
		zap.L().Info("query total task end.", zap.Int64("useTime", useTime))

	}()
	go func() {
		zap.L().Info("query task start...")
		queryTaskStartTime := time.Now()
		defer wait.Done()
		GetDB().Raw(querySql).Scan(resultBind)
		useTime := time.Now().Sub(queryTaskStartTime).Milliseconds()
		zap.L().Info("query task end.", zap.Int64("useTime", useTime))
	}()
	wait.Wait()
	zap.L().Info("page query finished", zap.Int64("totalUseTime", time.Now().Sub(startTime).Milliseconds()))

}
