package routes

import (
	"github.com/nilpanic/gin"

	"web-tpl/app/http/controllers/home"
)

func Reg(r *gin.Engine) {
	r.POST("/v1/user/add", home.Add)
}
