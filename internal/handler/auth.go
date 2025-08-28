package handler

import (
	"net/http"
	"shortURL/internal/model"
	"shortURL/internal/repo"
	"shortURL/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	userRepo *repo.UserRepo
}

func NewAuth(repo *repo.UserRepo) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{Email: req.Email, Password: string(hash)}
	if err := h.userRepo.Save(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct{ Email, Password string }
	if c.ShouldBindJSON(&req) != nil {
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := jwt.Generate(user.ID)
	c.JSON(200, gin.H{"token": token})
}
