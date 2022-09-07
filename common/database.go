// Package common
// @author： Boice
// @createTime：2022/7/22 15:45
package common

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/bns-engineering/td/common/util"
	"gitlab.com/bns-engineering/td/model/dto"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type Database struct {
	*gorm.DB
}

func newDatabase(config *Config, credential *Credential) *Database {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", credential.DB.Username, credential.DB.Password, credential.DB.Host, credential.DB.Port, credential.DB.Db)
	var err error
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields:            true,
		Logger:                 dbLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		util.CheckAndExit(err)
	}
	sqlDB, _ := db.DB()
	// Set the db connection configuration
	sqlDB.SetMaxOpenConns(credential.DB.MaxOpenConn) // set the max openning connection number
	sqlDB.SetMaxIdleConns(credential.DB.MaxIdleConn) // set the max idle connection number
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &Database{
		db,
	}

}

func Paginate(pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.Page <= 0 {
			pagination.Page = 1
		}
		if pagination.Perpage > 100 {
			pagination.Perpage = 100
		}
		offset := (pagination.Page - 1) * pagination.Perpage
		pagination.From = (pagination.Page-1)*pagination.Perpage + 1
		pagination.To = (pagination.Page-1)*pagination.Perpage + pagination.Perpage
		return db.Offset(int(offset)).Limit(int(pagination.Perpage))
	}
}

func FindPage(ctx context.Context, db *gorm.DB, pagination *dto.Pagination, resultBind interface{}) {
	L.Info(ctx, "page query start")
	startTime := time.Now()
	countSql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Count(&pagination.Total)
	})
	querySql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(Paginate(pagination)).Find(resultBind)
	})
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		L.Info(ctx, "query total task start...")
		totalTaskStartTime := time.Now()
		defer wait.Done()
		db.Raw(countSql).Scan(&pagination.Total)
		useTime := time.Now().Sub(totalTaskStartTime).Milliseconds()
		L.Info(ctx, "query total task end.", zap.Int64("useTime", useTime))
	}()
	go func() {
		L.Info(ctx, "query task start...")
		queryTaskStartTime := time.Now()
		defer wait.Done()
		db.Raw(querySql).Scan(resultBind)
		useTime := time.Now().Sub(queryTaskStartTime).Milliseconds()
		L.Info(ctx, "query task end.", zap.Int64("useTime", useTime))
	}()
	wait.Wait()
	// last page
	pagination.LastPage = decimal.NewFromInt(pagination.Total).Div(decimal.NewFromInt(int64(pagination.Perpage))).Ceil().IntPart()
	// to
	if pagination.To > pagination.Total {
		pagination.To = pagination.Total
	}
	L.Info(ctx, "page query finished", zap.Int64("totalUseTime", time.Now().Sub(startTime).Milliseconds()))

}
