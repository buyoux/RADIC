package dao

import (
	"RADIC/util"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"
)

var (
	search_mysql      *gorm.DB
	search_mysql_once sync.Once
	dblog             ormlog.Interface
)

func init() {
	// fout, err := os.OpenFile("log/mysql.log", os.O_CREATE, 0o644)
	// if err != nil {
	// 	util.LogRus.Panicf("open mysql.log failed: %s", err)
	// }
	dblog = ormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		ormlog.Config{
			SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      ormlog.Info,            // Log level, slient表示不输出日志
			Colorful:      true,                   // 彩色打印
		},
	)
}

func createMysqlDB(dbname, host, user, pass string, port int) *gorm.DB {
	// data source name : user_blog:wangs123@tcp(192.168.145.128:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)
	var err error
	// 启用PrepareStmt，SQL预编译，提高查询效率
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dblog, PrepareStmt: true})
	if err != nil {
		// panic(), 连接不到数据库直接panic，os.Exit(2)
		util.Log.Panicf("connect to mysql use dsn %s failed: %s", dsn, err)
	}
	// 设置数据库连接池参数，提高并发性能
	sqlDB, _ := db.DB()
	// 数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 连接池最大允许的空闲连接数，若sql任务需要执行的连接数大于20，超过20的会被连接池关闭
	sqlDB.SetMaxIdleConns(20)
	util.Log.Printf("connect to mysql db %s", dbname)
	return db
}

func GetSearchDBConnection() *gorm.DB {
	if search_mysql == nil {
		search_mysql_once.Do(func() {
			search_mysql = createMysqlDB("search", "192.168.145.128", "radic", "123123", 3306)
		})
	}
	return search_mysql
}
