package home

import (
	"github.com/gin-gonic/gin"
	"web-tpl/app/utils/rsp"
)

type Params struct {
	Page int `form:"page" binding:"required,gt=0,lt=1000"`
	Size int `form:"size,default=10" binding:"required,gt=0,lt=1000"`
}

func Index(ctx *gin.Context) {
	var prams Params

	if err := ctx.ShouldBind(&prams); err != nil {
		rsp.JSONErr(ctx, rsp.WithMsg(err.Error()))
		return
	}

	rsp.JSONOk(ctx, rsp.WithData(prams))
}
