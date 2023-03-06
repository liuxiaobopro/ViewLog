package resp

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	FailCode int = iota - 1
	SuccCode
	InterCode
	ParamCode
)

var codeMap = map[int]string{
	FailCode:  "fail",
	SuccCode:  "success",
	InterCode: "internal error",
	ParamCode: "param error",
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
		Msg:  "success",
		Data: data,
	}
}

func FailResp(args ...any) *Resp {
	code := FailCode
	msg := GetMsg(code)
	if len(args) > 0 {
		code = args[0].(int)
		msg = args[1].(string)
	}
	return &Resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
