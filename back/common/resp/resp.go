package resp

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	Fail int = iota - 1
	Succ
)

var codeMap = map[int]string{
	Fail: "fail",
}

func GetMsg(code int) string {
	return codeMap[code]
}

func SuccResp(data any) *Resp {
	return &Resp{
		Code: Succ,
		Msg:  "success",
		Data: data,
	}
}

func FailResp(args ...any) *Resp {
	code := Fail
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
