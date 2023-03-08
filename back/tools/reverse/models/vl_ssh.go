package models

import (
	"time"
)

type VlSsh struct {
	Id         int       `xorm:"not null pk autoincr INT"`
	Name       string    `xorm:"comment('主机名称') unique VARCHAR(50)"`
	IsActive   int       `xorm:"default 1 comment('活跃状态(1未连接 2已连接)') TINYINT"`
	Host       string    `xorm:"comment('主机地址') VARCHAR(50)"`
	Port       int       `xorm:"default 0 comment('端口号') INT"`
	Username   string    `xorm:"default '' comment('用户名') VARCHAR(50)"`
	Password   string    `xorm:"comment('密码') VARCHAR(200)"`
	CreateTime time.Time `xorm:"comment('创建时间') DATETIME"`
	UpdateTime time.Time `xorm:"comment('更新时间') DATETIME"`
}

func (m *VlSsh) TableComment() string {
	return "vl_ssh"
}
