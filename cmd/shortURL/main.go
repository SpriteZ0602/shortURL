package main

import (
	"log"
	"shortURL/internal/handler"
	"shortURL/internal/middleware"
	"shortURL/internal/model"
	"shortURL/internal/repo"
	"shortURL/internal/service"
	"shortURL/pkg/snowflake"
	"shortURL/pkg/trace"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 配置MySQL
	dsn := "root:123456@tcp(localhost:3306)/shorturl?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 迁移数据库, 传入 ShortURL 和 User 两个模型
	if err := db.AutoMigrate(&model.ShortURL{}, &model.User{}); err != nil {
		log.Fatal(err)
	}

	// 初始化雪花算法
	if err := snowflake.Init(); err != nil {
		log.Fatalf("snowflake init: %v", err)
	}

	// 创建数据库操作对象
	urlRepo := repo.New(db)

	// 创建业务逻辑对象
	svc := service.New(urlRepo, snowflake.Generate)

	// 创建回调对象
	h := handler.New(svc)

	// 创建用户操作对象
	userRepo := repo.NewUser(db)

	// 创建用户业务逻辑对象
	auth := handler.NewAuth(userRepo)

	// 创建路由
	r := gin.Default()
	r.POST("/api/v1/register", auth.Register)
	r.POST("/api/v1/login", auth.Login)
	r.GET("/:code", handler.NewRedirect(urlRepo))

	// 添加链路追踪中间件
	shutdown, _ := trace.Init("shorturl")
	defer shutdown()
	r.Use(middleware.Trace())

	// 添加JWT中间件的路由
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWT())
	authorized.POST("/shorten", h.Shorten)

	// 启动服务
	r.Run(":8080")
}
