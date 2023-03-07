package types

import (
	"fmt"
	"time"
)

type Time time.Time

//实现它的json序列化方法
func (th Time) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(th).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (th *Time) UnmarshalJSON(data []byte) error {
	local, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*th = Time(local)
	return err
}

type ModelExtends struct {
	CreateTime Time `xorm:"comment('创建时间') DATETIME created" json:"createTime"`
	UpdateTime Time `xorm:"comment('更新时间') DATETIME updated" json:"updateTime"`
	CreateUid  int  `xorm:"comment('创建人id') INT(11) default 0" json:"createUid"`
	UpdateUid  int  `xorm:"comment('更新人id') INT(11) default 0" json:"updateUid"`
}

type ModelCTExtends struct {
	CreateTime Time `xorm:"comment('创建时间') DATETIME created" json:"createTime"`
}

type ModelCUExtends struct {
	CreateTime Time `xorm:"comment('创建时间') DATETIME created" json:"createTime"`
	UpdateTime Time `xorm:"comment('更新时间') DATETIME updated" json:"updateTime"`
}
