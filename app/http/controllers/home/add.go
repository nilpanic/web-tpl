package home

import (
	"github.com/nilpanic/gin"
	"web-tpl/app/http/controllers/home/params"
	"web-tpl/app/utils/rsp"
)

func Add(ctx *gin.Context) {
	var p params.Add
	err := ctx.Valid(&p)
	if err != nil {
		//
		rsp.JSONErr(ctx, rsp.WithMsg(err.Error()))
		return
	}

	rsp.JSONOk(ctx, rsp.WithData(p))
}
