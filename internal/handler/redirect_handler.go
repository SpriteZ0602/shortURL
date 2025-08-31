// 短链接重定向回调

package handler

import (
	"context"
	"errors"
	"net/http"
	"shortURL/internal/repo"
	"shortURL/pkg/cache"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// NewRedirect 重定向
func NewRedirect(repo *repo.ShortURLRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		ctx := context.Background()

		// 1. 先查 Redis
		longURL, err := cache.RDB.Get(ctx, "short:"+code).Result()

		if errors.Is(err, redis.Nil) { // key 不存在
			// 2. 回源 MySQL
			su, err := repo.FindByCode(code)
			if err != nil || su == nil {
				c.String(http.StatusNotFound, "short code not found")
				return
			}
			// 回填缓存
			_ = cache.RDB.Set(ctx, "short:"+code, su.LongURL, 7*24*time.Hour)
			longURL = su.LongURL
		} else if err != nil {
			c.String(http.StatusInternalServerError, "cache error")
			return
		}

		// 3. 302 跳转
		c.Redirect(http.StatusFound, longURL)
	}
}
