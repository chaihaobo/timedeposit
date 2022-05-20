/*
 * @Author: Hugo
 * @Date: 2022-05-06 01:47:09
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-16 10:33:32
 */
package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Build the db handle for database
var _db *gorm.DB

// Initial gorm
func initDB() {
	// Do some initial for database
	// todo: move the init data to config file
	username := "hugo"                  //username
	password := "123456"                //password
	host := "172.17.0.2"                //host
	port := 3306                        //port
	Dbname := "time_deposit_eod_engine" //database name

	// Generate Mysql DSN
	// Mysql dsn format： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	// replace values which like {username}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	var err error
	// Connect Mysql, Get DB connection
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connect databse failed, error=" + err.Error())
	}

	sqlDB, _ := _db.DB()

	// Set the db connection configuration
	sqlDB.SetMaxOpenConns(100) // set the max openning connection number
	sqlDB.SetMaxIdleConns(20)  // set the max idle connection number
}

// Get db connection
func GetDB() *gorm.DB {
	if _db == nil {
		initDB()
	}
	return _db
}
