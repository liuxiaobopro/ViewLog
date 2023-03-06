package models

type VlDemo struct {
	Id      int64  `xorm:"pk autoincr BIGINT"`
	UsrName string `xorm:"not null comment('姓名') unique VARCHAR(25)"`
}

func (m *VlDemo) TableComment() string {
	return "vl_demo"
}
