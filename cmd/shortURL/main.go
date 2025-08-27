package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"shortURL/internal/handler"
	"shortURL/internal/model"
	"shortURL/internal/repo"
	"shortURL/internal/service"
	"shortURL/internal/util/shortID"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "root:123456@tcp(localhost:3306)/shorturl?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ShortURL{}); err != nil {
		log.Fatal(err)
	}

	repo := repo.New(db)
	svc := service.New(repo, shortID.New())
	h := handler.New(svc)

	r := gin.Default()
	r.POST("/api/v1/shorten", h.Shorten)
	r.Run(":8080")
}
