package global

import (
	"ViewLog/back/configs"

	"xorm.io/xorm"
)

var (
	Db   *xorm.Engine
	Conf *configs.Conf
)
