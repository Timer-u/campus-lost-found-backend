package main

import (
	"fmt"
	//"log"

	"campus-lost-found-backend/config"
	//"campus-lost-found-backend/model"
	"campus-lost-found-backend/router"
)

func main() {
	// 加载 config.yaml
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	//临时跳过数据库，优先跑通接口
	// model.InitDB()

	// // 程序退出时自动关闭数据库连接
	// sqlDB, err := model.DB.DB()
	// if err != nil {
	// 	log.Fatalf("获取数据库连接实例失败: %v", err)
	// }
	// defer sqlDB.Close()

	// 初始化 Gin 路由规则
	r := router.SetupRouter()

	// 拼接监听地址，例如 ":8080"
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)

	r.Run(addr)
}
