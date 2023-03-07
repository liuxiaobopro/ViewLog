package model

import "ViewLog/back/common/types"

type Demo struct {
	Id   int64
	Name string `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
}

type Ssh struct {
	Id       int    `xorm:"not null pk autoincr INT"`
	Name     string `xorm:"comment('主机名称') unique VARCHAR(50)"`
	Host     string `xorm:"comment('主机地址') VARCHAR(50)"`
	Port     int    `xorm:"default 0 comment('端口号') INT"`
	Username string `xorm:"default '' comment('用户名') VARCHAR(50)"`
	Password string `xorm:"default '' comment('密码') VARCHAR(50)"`

	types.ModelCUExtends `xorm:"extends"`
}
