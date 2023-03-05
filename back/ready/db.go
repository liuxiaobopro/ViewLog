package ready

import (
	"time"

	"ViewLog/back/global"
	"ViewLog/back/model"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func Db() {
	engine, err := xorm.NewEngine("mysql", "root:123456@127.0.0.1:3306/viewlog?charset=utf8mb4")
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)                     // 显示SQL语句
	engine.Logger().SetLevel(log.LOG_DEBUG)  // 设置日志级别
	engine.SetMaxIdleConns(2)                // 设置最大空闲连接数
	engine.SetMaxOpenConns(100)              // 设置最大连接数
	engine.SetConnMaxLifetime(4 * time.Hour) // 设置连接最大存活时间

	if err != nil {
		panic(err)
	}

	if err := engine.Sync(
		new(model.Demo),
	); err != nil {
		panic(err)
	}

	global.Db = engine
}
