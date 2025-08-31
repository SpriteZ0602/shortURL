package middleware

import (
	"context"
	"fmt"
	"net/http"
	"shortURL/pkg/cache"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

type QuotaMiddleware struct {
	dailyLimit int64
}

func NewQuota(dailyLimit int64) *QuotaMiddleware {
	return &QuotaMiddleware{dailyLimit: dailyLimit}
}

func (qm *QuotaMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := uid.(uint)
		key := fmt.Sprintf("quota:%d:%s", userID, time.Now().Format("2006-01-02"))
		cnt, _ := cache.RDB.Incr(ctx, key).Result()
		_ = cache.RDB.Expire(ctx, key, 24*time.Hour) // 改为24小时过期

		if cnt > qm.dailyLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "quota exceeded"})
			return
		}
		c.Next()
	}
}
