package service

import (
	"ViewLog/back/common/resp"
	"ViewLog/back/global"
	modelReq "ViewLog/back/model/req"
)

type apiService struct {
}

var ApiService = new(apiService)

// AddSsh 添加ssh
func (*apiService) AddSsh(req *modelReq.AddSshReq) *resp.Resp {
	var (
		sess = global.Db
	)
	sess.Where("name = ?", req.Name).Count()
	return resp.SuccResp(nil)
}
