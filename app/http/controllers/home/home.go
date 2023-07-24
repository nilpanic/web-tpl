package home

import (
	"github.com/gin-gonic/gin"
	"web-tpl/app"
)

func Index(ctx *gin.Context) {
	// app.log
	app.Log().Info("hello world")
	ctx.JSON(200, gin.H{
		"message": "hello world",
	})
}
