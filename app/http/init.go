package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"web-tpl/app"
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

	//_ = r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5173"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 其它配置
	routes.Reg(r)

	return r.Run(app.Config.HTTP.Listen) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
