package model

type Demo struct {
	Id   int64
	Name string `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
}
