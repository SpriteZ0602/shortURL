// 注册登录回调

package handler

import (
	"net/http"
	"shortURL/internal/model"
	"shortURL/internal/repo"
	"shortURL/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterReq 注册请求
type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginReq 登录请求结构体
type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo *repo.UserRepo
}

// NewAuth 创建一个AuthHandler
func NewAuth(repo *repo.UserRepo) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	// 绑定请求体
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成密码哈希
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{Email: req.Email, Password: string(hash)}

	// 保存用户
	if err := h.userRepo.Save(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	// 绑定请求体
	var req struct{ Email, Password string }

	// 解析请求体
	if c.ShouldBindJSON(&req) != nil {
		return
	}

	// 获取用户
	user, err := h.userRepo.FindByEmail(req.Email)

	// 验证用户
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// 生成令牌
	token, _ := jwt.Generate(user.ID)
	c.JSON(200, gin.H{"token": token})
}
