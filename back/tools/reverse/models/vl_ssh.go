package models

type VlSsh struct {
	Id int `xorm:"not null pk autoincr INT"`
}

func (m *VlSsh) TableComment() string {
	return "vl_ssh"
}
