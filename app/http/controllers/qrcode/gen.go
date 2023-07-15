package qrcode

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

type params struct {
	URL string `form:"url" binding:"required,url"`
}

func Gen(ctx *gin.Context) {
	// 验证请求参数
	var p params
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 生成二维码
	var png []byte
	png, err := qrcode.Encode(p.URL, qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"url": "data:image/png;base64," + base64.StdEncoding.EncodeToString(png),
		},
	})
}
