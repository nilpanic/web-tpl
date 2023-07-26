package routes

import (
	"github.com/gin-gonic/gin"

	"web-tpl/app/http/controllers/home"
)

func Reg(r *gin.Engine) {
	r.GET("/v1/user", home.Index)
}
