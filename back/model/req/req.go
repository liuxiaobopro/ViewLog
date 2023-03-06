package req

type LogIndexReadFileReq struct {
	Path string `form:"path" json:"path" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
	Page int    `form:"page" json:"page" binding:"required"`
}

type AddSshReq struct {
	Name string `form:"name" json:"name" binding:"required"`
}
