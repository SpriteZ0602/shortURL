package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("shorturl")
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
