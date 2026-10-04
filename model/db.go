package model

import (
	"fmt"
	"log"
	"time"

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

	// 连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接实例失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)                  // 最大同时打开的连接数
	sqlDB.SetMaxIdleConns(5)                   // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 单个连接最长存活时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 空闲连接最长存活时间

	err = DB.AutoMigrate(&User{}, &Item{}, &Claim{}, &Announcement{})
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	log.Println("数据库表结构初始化/迁移成功")
}
