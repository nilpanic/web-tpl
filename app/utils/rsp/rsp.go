package rsp

import (
	"github.com/nilpanic/gin"
	"net/http"
)

// 函数编程设计
// {
// 	"code": 10000,
// 	"msg": "success",
// 	"data": {}
// }

type JSONVal struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type JsonOption func(*JSONVal)

func WithCode(code int) JsonOption {
	return func(v *JSONVal) {
		v.Code = code
	}
}

func WithMsg(msg string) JsonOption {
	return func(v *JSONVal) {
		v.Msg = msg
	}
}

func WithData(data any) JsonOption {
	return func(v *JSONVal) {
		v.Data = data
	}
}

func JSONOk(ctx *gin.Context, opts ...JsonOption) {
	v := JSONVal{
		Code: 0,
		Msg:  "success",
	}

	for _, opt := range opts {
		opt(&v)
	}

	ctx.JSON(http.StatusOK, v)
}

func JSONErr(ctx *gin.Context, opts ...JsonOption) {
	v := JSONVal{
		Code: 1,
		Msg:  "error",
	}

	for _, opt := range opts {
		opt(&v)
	}

	ctx.JSON(http.StatusOK, v)
}
