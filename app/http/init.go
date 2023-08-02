package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/nilpanic/gin"
	"github.com/nilpanic/gin/binding"
	"web-tpl/app"
	"web-tpl/app/core/valid"
	"web-tpl/app/http/middleware/logger"
	"web-tpl/app/http/routes"
)

func NewServer() error {
	// 启动 server
	r := gin.New()
	if app.Config.WebServerLog.Enable {
		r.Use(logger.New(app.Config.WebServerLog, app.Config.Env, app.Config.HomeDir))
	}
	r.Use(gin.Recovery())

	// register validator func
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", valid.Mobile)
	}

	// 其它配置
	routes.Reg(r)

	return r.Run(app.Config.HTTP.Listen) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
