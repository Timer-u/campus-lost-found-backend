package model

import (
	"fmt"
	"log"

	"campus-lost-found-backend/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GlobalConfig.Database
	// 拼接 MySQL 数据源名称 (DSN)
	// 格式：user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBname,
	)

	var err error
	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")

	err = DB.AutoMigrate(&User{}, &Item{})
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	log.Println("数据库表结构初始化/迁移成功")
}
