// 生成短链接回调

package handler

import (
	"log"
	"net/http"
	"shortURL/internal/middleware"
	"shortURL/internal/service"

	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
)

type ShortHandler struct{ svc *service.ShortService }

func New(s *service.ShortService) *ShortHandler { return &ShortHandler{svc: s} }

func (h *ShortHandler) Shorten(c *gin.Context) {

	ctx, span := otel.Tracer("shorturl").Start(c.Request.Context(), "handler.Shorten")
	defer span.End()

	var req struct {
		URL string `json:"url" binding:"required"`
	}

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否需要进行风险检查（通过中间件设置的标记）
	if needCheck, exists := c.Get("needRiskCheck"); exists && needCheck.(bool) {
		// 进行URL风险检查
		if middleware.CheckURLRisk(req.URL) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "blocked url"})
			return
		}
	}

	// 生成短链接
	code, err := h.svc.Shorten(ctx, req.URL)
	if err != nil {
		log.Println("Shorten error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_code": code})
}
