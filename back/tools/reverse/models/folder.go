package models

type Folder struct {
	Id    int    `xorm:"not null pk autoincr INT"`
	SshId int    `xorm:"default 0 comment('ssh自增id') INT"`
	Name  string `xorm:"comment('名称') unique VARCHAR(50)"`
	Path  string `xorm:"default '' comment('文件夹全路径') VARCHAR(200)"`
}

func (m *Folder) TableComment() string {
	return "folder"
}
