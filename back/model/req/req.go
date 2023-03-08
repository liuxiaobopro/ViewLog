package req

type InstallReq struct {
	Host     string `form:"host" json:"host" binding:"required"`
	Port     int    `form:"port" json:"port" binding:"required"`
	Dbname   string `form:"dbname" json:"dbname" binding:"required"`
	User     string `form:"user" json:"user"`
	Password string `form:"password" json:"password"`
	Charset  string `form:"charset" json:"charset" binding:"required"`
}

type LogIndexReadFileReq struct {
	Path string `form:"path" json:"path" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
	Page int    `form:"page" json:"page" binding:"required"`
}

type AddSshReq struct {
	Name     string `form:"name" json:"name" binding:"required"`
	IsActive int    `form:"isActive" json:"isActive"`
	Host     string `form:"host" json:"host" binding:"required"`
	Port     int    `form:"port" json:"port" binding:"required"`
	Username string `form:"user" json:"username"`
	Password string `form:"password" json:"password"`
}

type DelSshReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

type UpdateSshReq struct {
	Id       int    `form:"id" json:"id" binding:"required"`
	Name     string `form:"name" json:"name"`
	IsActive int    `form:"isActive" json:"isActive"`
	Host     string `form:"host" json:"host"`
	Port     int    `form:"port" json:"port"`
	Username string `form:"user" json:"username"`
	Password string `form:"password" json:"password"`
}

type DetailSshReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

type ListSshReq struct {
	Page  int `form:"page" json:"page" binding:"required"`
	Limit int `form:"limit" json:"limit" binding:"required"`
}
