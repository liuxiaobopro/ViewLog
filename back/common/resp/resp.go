package resp

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ListResult struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

const (
	FailCode int = iota - 1
	SuccCode
	InterCode
	ParamCode
)

var codeMap = map[int]string{
	FailCode:  "操作失败",
	SuccCode:  "操作成功",
	InterCode: "内部错误",
	ParamCode: "参数错误",
}

func GetMsg(code int) string {
	return codeMap[code]
}

var (
	Fail  = &Resp{Code: FailCode, Msg: GetMsg(FailCode), Data: nil}
	Succ  = &Resp{Code: SuccCode, Msg: GetMsg(SuccCode), Data: nil}
	Inter = &Resp{Code: InterCode, Msg: GetMsg(InterCode), Data: nil}
	Param = &Resp{Code: ParamCode, Msg: GetMsg(ParamCode), Data: nil}
)

func SuccResp(data any) *Resp {
	return &Resp{
		Code: SuccCode,
		Msg:  GetMsg(SuccCode),
		Data: data,
	}
}

func FailResp(args ...any) *Resp {
	code := FailCode
	msg := GetMsg(code)
	if len(args) > 0 {
		code = args[0].(int)   // 传入的code
		msg = args[1].(string) // 传入的msg
	}
	return &Resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
