package middleware

import (
	"context"
	"shortURL/pkg/cache"

	"github.com/gin-gonic/gin"
)

// RiskCheck 重写风险检查中间件，避免解析请求体导致的EOF问题
// 改为在handler中进行URL检查，或者如果必须在此处检查，使用更安全的方式
func RiskCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不直接解析请求体，而是在请求上下文中设置一个标记
		// 实际的URL检查将在handler中进行
		c.Set("needRiskCheck", true)
		c.Next()
	}
}

// CheckURLRisk 提供一个函数供handler调用，进行实际的URL风险检查
func CheckURLRisk(url string) bool {
	ctx := context.Background()
	return cache.RDB.SIsMember(ctx, "blacklist", url).Val()
}
