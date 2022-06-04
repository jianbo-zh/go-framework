package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LimitBodySize(c *gin.Context) {

	// 32MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32<<20)

	c.Next()
}
