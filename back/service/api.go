package service

import (
	"ViewLog/back/common/resp"
	modelReq "ViewLog/back/model/req"
)

type apiService struct {
}

var ApiService = new(apiService)

// AddSsh 添加ssh
func (*apiService) AddSsh(req *modelReq.AddSshReq) *resp.Resp {
	return resp.SuccResp(nil)
}
