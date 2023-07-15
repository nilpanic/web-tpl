package home

import (
	"github.com/gin-gonic/gin"

	"web-tpl/app/http/models"
)

func Index(ctx *gin.Context) {
	var rel []models.User


	ctx.JSON(200, gin.H{
		"code": 0,
		"data": rel,
	})
}
