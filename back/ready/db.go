package ready

import (
	"fmt"
	"time"

	"ViewLog/back/global"
	"ViewLog/back/model"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func Db() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		global.Conf.DB.User,
		global.Conf.DB.Password,
		global.Conf.DB.Host,
		global.Conf.DB.Port,
		global.Conf.DB.Name,
		global.Conf.DB.Charset)
	engine, err := xorm.NewEngine("mysql", dns)
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
