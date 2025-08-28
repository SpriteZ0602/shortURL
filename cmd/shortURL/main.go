package main

import (
	"log"
	"shortURL/internal/handler"
	"shortURL/internal/middleware"
	"shortURL/internal/model"
	"shortURL/internal/repo"
	userRepoPkg "shortURL/internal/repo"
	"shortURL/internal/service"
	"shortURL/internal/util/shortID"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(localhost:3306)/shorturl?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ShortURL{}, &model.User{}); err != nil {
		log.Fatal(err)
	}

	repo := repo.New(db)
	svc := service.New(repo, shortID.New())
	h := handler.New(svc)

	userRepo := userRepoPkg.NewUser(db)
	auth := handler.NewAuth(userRepo)

	r := gin.Default()
	r.POST("/api/v1/register", auth.Register)
	r.POST("/api/v1/login", auth.Login)
	r.GET("/:code", handler.NewRedirect(repo))

	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWT())
	authorized.POST("/shorten", h.Shorten)
	r.Run(":8080")
}
