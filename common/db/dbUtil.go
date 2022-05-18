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

// 定义全局的db对象，我们执行数据库操作主要通过他实现。
var _db *gorm.DB

// 包初始化函数，golang 特性，每个包初始化的时候会自动执行 init 函数，这里用来初始化 gorm。
func initDB() {
	// ...忽略dsn配置，请参考上面例子...
	// 配置 MySQL 连接参数
	username := "hugo"                  //账号
	password := "123456"                //密码
	host := "172.17.0.2"                //数据库地址，可以是Ip或者域名
	port := 3306                        //数据库端口
	Dbname := "time_deposit_eod_engine" //数据库名

	// 通过前面的数据库参数，拼接 Mysql DSN，其实就是数据库连接串（数据源名称）
	// Mysql dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	// 类似{username}使用花括号包着的名字都是需要替换的参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	// 声明 err 变量，下面不能使用 := 赋值运算符，否则_db变量会当成局部变量，导致外部无法访问 _db 变量
	var err error
	// 连接 Mysql, 获得DB类型实例，用于后面的数据库读写操作。
	fmt.Println(dsn)
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	sqlDB, _ := _db.DB()

	// 设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
}

// 获取 gorm db 对象，其他包需要执行数据库查询的时候，只要通过 tools.getDB() 获取 db 对象即可。
// 不用担心协程并发使用同样的db对象会共用同一个连接，db 对象在调用他的方法的时候会从数据库连接池中获取新的连接
func GetDB() *gorm.DB {
	if _db == nil {
		initDB()
	}
	return _db
}
