package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"llm-security-bas/internal/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v\n%s", err, debug.Stack())
				c.JSON(http.StatusOK, response.Response{
					Code: response.CodeInternalError,
					Msg:  "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
