package home

import (
	"github.com/nilpanic/gin"
	"web-tpl/app/http/controllers/home/params"
	"web-tpl/app/utils/rsp"
)

func Index(ctx *gin.Context) {
	var prams params.Add

	if err := ctx.ShouldBind(&prams); err == nil {
		rsp.JSONOk(ctx, rsp.WithData(prams))
	} else {
		rsp.JSONErr(ctx, rsp.WithMsg(err.Error()))
	}

}
