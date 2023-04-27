package ginhook

import "fmt"

var (
	CODE_OK        = 200 //接口正常
	CODE_ERROR     = 201 //常规错误
	CODE_REDIRECT  = 302 //需要重新登录，调用重新授权也可能成功
	CODE_NOT_FOUND = 404 //程序报错
)

type Exception struct {
	Code int
	Msg  any
	Data any
}

func (e *Exception) Error() string {
	return fmt.Sprint(e.Msg)
}

func ThrowError(msg string) {
	panic(Exception{
		Msg:  msg,
		Code: CODE_ERROR,
		Data: map[string]interface{}{},
	})
}

func ThrowErrorCode(msg string, code int) {
	panic(Exception{
		Msg:  msg,
		Code: code,
		Data: map[string]interface{}{},
	})
}
func ThrowErrorCodeData(msg string, code int, data any) {
	panic(Exception{
		Msg:  msg,
		Code: code,
		Data: data,
	})
}

func ThrowError302(msg string) {
	panic(Exception{
		Msg:  msg,
		Code: CODE_REDIRECT,
		Data: map[string]interface{}{},
	})
}

func Error(msg any) error {
	return &Exception{
		Msg:  msg,
		Code: CODE_ERROR,
		Data: map[string]interface{}{},
	}
}

func ErrorData(msg any, data any) error {
	return &Exception{
		Msg:  msg,
		Code: CODE_ERROR,
		Data: data,
	}
}

func ErrorCode(msg any, code int) error {
	return &Exception{
		Msg:  msg,
		Code: code,
		Data: map[string]interface{}{},
	}
}
func Error302(msg any) error {
	return &Exception{
		Msg:  msg,
		Code: CODE_REDIRECT,
		Data: map[string]interface{}{},
	}
}

func Try(c func(data Exception)) {
	if e := recover(); e != nil {
		if v, ok := e.(Exception); ok {
			c(v)
		} else {
			//代码异常
			c(Exception{
				Code: CODE_NOT_FOUND,
				Msg:  e,
				Data: map[string]interface{}{},
			})
		}
	}
}
