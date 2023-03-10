package model

import "ViewLog/back/common/types"

type Demo struct {
	Id   int64
	Name string `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
}

type Ssh struct {
	Id       int    `xorm:"not null pk autoincr INT" json:"id"`
	Name     string `xorm:"comment('主机名称') unique VARCHAR(50)" json:"name"`
	IsActive int    `xorm:"default 1 comment('活跃状态(1未连接 2已连接)') TINYINT" json:"isActive"`
	Host     string `xorm:"comment('主机地址') VARCHAR(50)" json:"host"`
	Port     int    `xorm:"default 0 comment('端口号') INT" json:"port"`
	Username string `xorm:"default '' comment('用户名') VARCHAR(50)" json:"username"`
	Password string `xorm:"default '' comment('密码') VARCHAR(50)" json:"-"`

	types.ModelCUExtends `xorm:"extends"`
}

type Folder struct {
	Id    int    `xorm:"not null pk autoincr INT"`
	SshId int    `xorm:"default 0 comment('ssh自增id') INT"`
	Name  string `xorm:"comment('名称') unique VARCHAR(50)"`
	Path  string `xorm:"default '' comment('文件夹全路径') VARCHAR(200)"`

	types.ModelCUExtends `xorm:"extends"`
}
