package main

import (
	"fmt"

	"campus-lost-found-backend/config"
	"campus-lost-found-backend/model"
	"campus-lost-found-backend/router"
)

func main() {
	// 加载 config.yaml
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	model.InitDB()

	// 初始化 Gin 路由规则
	r := router.SetupRouter()

	// 拼接监听地址，例如 ":8080"
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	r.Run(addr)
}
