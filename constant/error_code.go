package constant

import "fmt"

type ErrorCode struct {
	Code int64
	Msg  string
}

func makeErrorCode(code int64, msg string) *ErrorCode {
	return &ErrorCode{
		Code: code,
		Msg:  msg,
	}
}

func BuildErrCode(detail interface{}, errCode *ErrorCode) *ErrorCode {
	return makeErrorCode(errCode.Code, fmt.Sprintf(errCode.Msg, detail))
}

var (
	RetSuccess              = makeErrorCode(0, "成功")
	RetParamsErr            = makeErrorCode(1, "参数非法:%v")
	RetNoSuchUserErr        = makeErrorCode(2, "用户名不存在或者密码错误:%v")
	RetDuplicateUsernameErr = makeErrorCode(3, "用户名重复，请修改用户名:%v")
	RetWrongGenderErr       = makeErrorCode(4, "性别不正确，请修改性别:%v")
	RetInvalidTokenErr      = makeErrorCode(5, "传入的token失效:%v")
	RetWriteRepoErr         = makeErrorCode(100, "写库错误:%v")
	RetReadRepoErr          = makeErrorCode(101, "查库错误:%v")
	RetSysErr               = makeErrorCode(999, "系统未知错误:%v")
)
