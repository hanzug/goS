package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			// 检查是否有panic
			if r := recover(); r != nil {
				context.JSON(500, gin.H{
					"code": 500,
					"msg":  fmt.Sprintf("%s", r),
				})
				context.Abort()
			}
		}()
		context.Next()
	}
}
