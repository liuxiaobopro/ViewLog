package ready

import (
	"ViewLog/back/global"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

func Db() {
	engine, err := xorm.NewEngine("sqlite3", "./test.db")

	if err != nil {
		panic(err)
	}

	global.Db = engine
}
