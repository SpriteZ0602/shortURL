// 生成短链接回调

package handler

import (
	"log"
	"net/http"
	"shortURL/internal/service"

	"github.com/gin-gonic/gin"
)

type ShortHandler struct{ svc *service.ShortService }

func New(s *service.ShortService) *ShortHandler { return &ShortHandler{svc: s} }

func (h *ShortHandler) Shorten(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成短链接
	code, err := h.svc.Shorten(req.URL)
	if err != nil {
		log.Println("Shorten error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_code": code})
}
