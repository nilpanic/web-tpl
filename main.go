package main

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"

	"web-tpl/app"
	"web-tpl/app/http"
)

func main() {
	// 解析项目根目录参数
	homePath := flag.String("prjHome", "", "项目的根目录路径")
	flag.Parse()

	if *homePath == "" {
		_, f, _, ok := runtime.Caller(0)
		if !ok {
			panic("尝试获取文件路径失败！")
		}

		*homePath = filepath.Dir(f)
	}

	err := app.Init(*homePath)
	if err != nil {
		panic(err)
	}

	err = http.NewServer()
	if err != nil {
		log.Fatal(err)
	}
}
